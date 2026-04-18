package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/config"
	"vaultpull/vault"
)

func init() {
	var dir string
	var snapshotFile string

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Save or restore a snapshot of secrets from Vault",
	}

	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Save current Vault secrets to a snapshot file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
			if err != nil {
				return err
			}
			secrets, err := vault.FetchSecrets(client, cfg.SecretPath)
			if err != nil {
				return err
			}
			out, err := vault.SaveSnapshot(dir, cfg.SecretPath, secrets)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Snapshot saved: %s\n", out)
			return nil
		},
	}
	saveCmd.Flags().StringVar(&dir, "dir", ".snapshots", "Directory to store snapshots")

	restoreCmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore secrets from a snapshot file to a .env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if snapshotFile == "" {
				return fmt.Errorf("--file is required")
			}
			s, err := vault.LoadSnapshot(snapshotFile)
			if err != nil {
				return err
			}
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if err := vault.WriteEnvFile(cfg.OutputFile, s.Secrets); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Restored %d secrets to %s\n", len(s.Secrets), cfg.OutputFile)
			return nil
		},
	}
	restoreCmd.Flags().StringVar(&snapshotFile, "file", "", "Snapshot file to restore from")

	snapshotCmd.AddCommand(saveCmd, restoreCmd)
	rootCmd.AddCommand(snapshotCmd)
}
