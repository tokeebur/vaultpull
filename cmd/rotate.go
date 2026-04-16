package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/vaultpull/config"
	"github.com/user/vaultpull/vault"
)

func init() {
	rotateCmd := &cobra.Command{
		Use:   "rotate",
		Short: "Rotate secrets: backup current .env and write fresh values from Vault",
		RunE:  runRotate,
	}
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.Token, cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	result, err := vault.RotateSecrets(client, cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	fmt.Fprintf(os.Stdout, "✔ Rotated secrets into %s\n", result.OutputFile)
	if result.BackupFile != "" {
		fmt.Fprintf(os.Stdout, "  backup: %s\n", result.BackupFile)
	}
	fmt.Fprintf(os.Stdout, "  added=%d removed=%d changed=%d\n",
		len(result.Diff.Added),
		len(result.Diff.Removed),
		len(result.Diff.Changed),
	)
	return nil
}
