package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourorg/vaultpull/config"
	"github.com/yourorg/vaultpull/vault"
)

var (
	snapshotFile   string
	snapshotList   bool
	snapshotRestore bool
)

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Save, list, or restore snapshots of Vault secrets",
		RunE:  runSnapshot,
	}

	snapshotCmd.Flags().StringVarP(&snapshotFile, "file", "f", "", "Snapshot file path (for save/restore)")
	snapshotCmd.Flags().BoolVarP(&snapshotList, "list", "l", false, "List contents of a snapshot file")
	snapshotCmd.Flags().BoolVarP(&snapshotRestore, "restore", "r", false, "Restore secrets from a snapshot file to the output .env file")

	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	// List mode: print snapshot contents
	if snapshotList {
		if snapshotFile == "" {
			return fmt.Errorf("--file is required with --list")
		}
		snap, err := vault.LoadSnapshot(snapshotFile)
		if err != nil {
			return fmt.Errorf("failed to load snapshot: %w", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "KEY\tVALUE")
		for k, v := range snap {
			fmt.Fprintf(w, "%s\t%s\n", k, v)
		}
		w.Flush()
		return nil
	}

	// Restore mode: write snapshot back to .env file
	if snapshotRestore {
		if snapshotFile == "" {
			return fmt.Errorf("--file is required with --restore")
		}
		snap, err := vault.LoadSnapshot(snapshotFile)
		if err != nil {
			return fmt.Errorf("failed to load snapshot: %w", err)
		}
		if err := vault.WriteEnvFile(cfg.OutputFile, snap); err != nil {
			return fmt.Errorf("failed to restore snapshot to %s: %w", cfg.OutputFile, err)
		}
		fmt.Printf("Restored %d keys from snapshot to %s\n", len(snap), cfg.OutputFile)
		return nil
	}

	// Default: fetch secrets and save a new snapshot
	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client error: %w", err)
	}

	secrets, err := client.FetchSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("failed to fetch secrets: %w", err)
	}

	dest := snapshotFile
	if dest == "" {
		dest = fmt.Sprintf(".snapshot-%s.json", time.Now().Format("20060102-150405"))
	}

	if err := vault.SaveSnapshot(dest, secrets); err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	data, _ := json.Marshal(secrets)
	_ = data
	fmt.Printf("Snapshot saved: %s (%d keys)\n", dest, len(secrets))
	return nil
}
