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
	var inPlace bool
	var outputFile string

	cmd := &cobra.Command{
		Use:   "placeholder",
		Short: "Replace secret values with placeholder text for safe sharing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPlaceholder(rulesFlag, inPlace, outputFile)
		},
	}

	cmd.Flags().StringSliceVar(&rulesFlag, "rule", nil, "Extra rules in PATTERN=PLACEHOLDER format (e.g. '*_CERT=<cert>'")
	cmd.Flags().BoolVar(&inPlace, "in-place", false, "Overwrite the output file with placeholdered values")
	cmd.Flags().StringVar(&outputFile, "output", "", "Output file (defaults to config output_file)")

	rootCmd.AddCommand(cmd)
}

func runPlaceholder(rulesFlag []string, inPlace bool, outputFile string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	target := cfg.OutputFile
	if outputFile != "" {
		target = outputFile
	}

	secrets, err := vault.ParseEnvFile(target)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	rules := vault.DefaultPlaceholderRules()
	for _, r := range rulesFlag {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid rule %q: expected PATTERN=PLACEHOLDER", r)
		}
		rules = append(rules, vault.PlaceholderRule{Pattern: parts[0], Placeholder: parts[1]})
	}

	result, err := vault.ApplyPlaceholders(secrets, rules)
	if err != nil {
		return fmt.Errorf("applying placeholders: %w", err)
	}

	changed := vault.ListPlaceholderKeys(secrets, result)
	if len(changed) == 0 {
		fmt.Println("No keys matched placeholder rules.")
		return nil
	}

	dest := os.Stdout
	if inPlace {
		f, ferr := os.Create(target)
		if ferr != nil {
			return fmt.Errorf("opening output file: %w", ferr)
		}
		defer f.Close()
		dest = f
	}

	for _, k := range changed {
		fmt.Fprintf(dest, "%s=%s\n", k, result[k])
	}
	fmt.Fprintf(os.Stderr, "Placeholdered %d key(s): %s\n", len(changed), strings.Join(changed, ", "))
	return nil
}
