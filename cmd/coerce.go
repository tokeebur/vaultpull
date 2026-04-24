package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/config"
	"vaultpull/vault"
)

func init() {
	var rulesFlag []string
	var outputFile string
	var dryRun bool

	coerceCmd := &cobra.Command{
		Use:   "coerce",
		Short: "Apply coerce rules (trim, upper, lower) to secret values",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCoerce(rulesFlag, outputFile, dryRun)
		},
	}

	coerceCmd.Flags().StringArrayVarP(&rulesFlag, "rule", "r", nil, "Coerce rule in KEY:ACTION format (e.g. API_KEY:trim)")
	coerceCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output .env file (defaults to config output file)")
	coerceCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print report without writing changes")

	rootCmd.AddCommand(coerceCmd)
}

func runCoerce(rulesFlag []string, outputFile string, dryRun bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	if outputFile == "" {
		outputFile = cfg.OutputFile
	}

	existing, err := vault.ParseEnvFile(outputFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read env file: %w", err)
	}

	var rules []vault.CoerceRule
	for _, r := range rulesFlag {
		parts := strings.SplitN(r, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid rule %q: expected KEY:ACTION", r)
		}
		rules = append(rules, vault.CoerceRule{Key: parts[0], Action: parts[1]})
	}

	result, report, err := vault.CoerceSecrets(existing, rules)
	if err != nil {
		return fmt.Errorf("coerce: %w", err)
	}

	fmt.Println(vault.FormatCoerceReport(report))

	if dryRun {
		fmt.Println("(dry-run) no changes written")
		return nil
	}

	if err := vault.WriteEnvFile(outputFile, result); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}

	fmt.Printf("wrote coerced secrets to %s\n", outputFile)
	return nil
}
