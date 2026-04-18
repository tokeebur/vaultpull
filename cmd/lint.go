package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/config"
	"vaultpull/vault"
)

func init() {
	var lintCmd = &cobra.Command{
		Use:   "lint",
		Short: "Lint secrets fetched from Vault for common issues",
		RunE:  runLint,
	}
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
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

	rules := vault.DefaultRules()
	results := vault.LintSecrets(secrets, rules)

	if len(results) == 0 {
		fmt.Println("✔ No lint violations found.")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Lint violations found (%d):\n", len(results))
	for _, r := range results {
		fmt.Fprintf(os.Stderr, "  [%s] %s\n", r.Rule, r.Message)
	}
	return fmt.Errorf("lint failed with %d violation(s)", len(results))
}
