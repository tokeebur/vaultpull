package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/config"
	"vaultpull/vault"
)

func init() {
	var rulesFile string
	var outputFile string

	deprecateCmd := &cobra.Command{
		Use:   "deprecate",
		Short: "Check secrets for deprecated keys",
		Long:  "Scans the local .env file against a deprecation rules file and reports any deprecated keys still in use.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeprecate(rulesFile, outputFile)
		},
	}

	deprecateCmd.Flags().StringVarP(&rulesFile, "rules", "r", ".deprecations", "Path to deprecation rules file")
	deprecateCmd.Flags().StringVarP(&outputFile, "file", "f", ".env", "Path to .env file to check")

	rootCmd.AddCommand(deprecateCmd)
}

func runDeprecate(rulesFile, envFile string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	_ = cfg

	rules, err := vault.LoadDeprecationFile(rulesFile)
	if err != nil {
		return fmt.Errorf("loading deprecation rules: %w", err)
	}
	if len(rules) == 0 {
		fmt.Println("No deprecation rules defined.")
		return nil
	}

	secrets, err := vault.ParseEnvFile(envFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading env file: %w", err)
	}

	report := vault.CheckDeprecations(secrets, rules)
	fmt.Print(vault.FormatDeprecationReport(report))

	if len(report.Deprecated) > 0 {
		os.Exit(1)
	}
	return nil
}
