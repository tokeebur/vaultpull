package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/vaultpull/config"
	"github.com/user/vaultpull/vault"
)

var (
	secretPath string
	outputFile string
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull secrets from Vault and write to a .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(secretPath, outputFile)
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}

		client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
		if err != nil {
			return fmt.Errorf("vault client error: %w", err)
		}

		secrets, err := client.GetSecrets(cfg.SecretPath)
		if err != nil {
			return fmt.Errorf("failed to fetch secrets: %w", err)
		}

		if err := vault.WriteEnvFile(cfg.OutputFile, secrets); err != nil {
			return fmt.Errorf("failed to write env file: %w", err)
		}

		fmt.Fprintf(os.Stdout, "✓ Wrote %d secrets to %s\n", len(secrets), cfg.OutputFile)
		return nil
	},
}

func init() {
	pullCmd.Flags().StringVarP(&secretPath, "path", "p", "", "Vault secret path (required)")
	pullCmd.Flags().StringVarP(&outputFile, "output", "o", ".env", "Output .env file path")
	pullCmd.MarkFlagRequired("path")
}
