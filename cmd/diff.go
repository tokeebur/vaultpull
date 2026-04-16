package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/config"
	"vaultpull/vault"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show differences between Vault secrets and local .env file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}

		client, err := vault.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("vault client error: %w", err)
		}

		secrets, err := client.FetchSecrets(cfg.SecretPath)
		if err != nil {
			return fmt.Errorf("fetch error: %w", err)
		}

		existing, err := vault.ParseEnvFile(cfg.OutputFile)
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		diff := vault.ComputeDiff(existing, secrets)

		if len(diff.Added)+len(diff.Changed)+len(diff.Removed) == 0 {
			fmt.Println("No changes detected.")
			return nil
		}

		for k, v := range diff.Added {
			fmt.Fprintf(os.Stdout, "+ %s=%s\n", k, v)
		}
		for k, v := range diff.Changed {
			fmt.Fprintf(os.Stdout, "~ %s=%s\n", k, v)
		}
		for k := range diff.Removed {
			fmt.Fprintf(os.Stdout, "- %s\n", k)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
