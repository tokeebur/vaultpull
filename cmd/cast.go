package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/config"
	"vaultpull/vault"
)

var castRulesFile string

func init() {
	castCmd := &cobra.Command{
		Use:   "cast",
		Short: "Cast secret values to typed representations and display results",
		RunE:  runCast,
	}
	castCmd.Flags().StringVar(&castRulesFile, "rules", "cast-rules.env", "File with KEY=type cast rules")
	rootCmd.AddCommand(castCmd)
}

func runCast(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	secrets, err := client.FetchSecrets(cfg.SecretPath)
	if err != nil {
		return err
	}

	rawRules, err := vault.ParseEnvFile(castRulesFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading cast rules: %w", err)
	}

	var rules []vault.CastRule
	for k, v := range rawRules {
		rules = append(rules, vault.CastRule{Pattern: k, Type: vault.CastType(v)})
	}

	results, err := vault.CastSecrets(secrets, rules)
	if err != nil {
		return err
	}

	for _, r := range results {
		fmt.Printf("%-30s %-8s %v\n", r.Key, r.Type, r.Cast)
	}
	return nil
}
