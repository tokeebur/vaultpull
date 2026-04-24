package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/config"
	"github.com/yourusername/vaultpull/vault"
)

func init() {
	var scopeName string
	var scopeFile string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "scope",
		Short: "Apply a named scope filter to secrets before writing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScope(scopeName, scopeFile, outputFile)
		},
	}
	cmd.Flags().StringVar(&scopeName, "name", "", "Scope name to apply (required)")
	cmd.Flags().StringVar(&scopeFile, "scope-file", ".vaultscopes", "Path to scope definition file")
	cmd.Flags().StringVar(&outputFile, "output", "", "Output .env file (defaults to config output_file)")
	_ = cmd.MarkFlagRequired("name")
	rootCmd.AddCommand(cmd)
}

func runScope(name, scopeFile, outputFile string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	if outputFile == "" {
		outputFile = cfg.OutputFile
	}

	scopes, err := vault.LoadScopeFile(scopeFile)
	if err != nil {
		return fmt.Errorf("load scope file: %w", err)
	}

	var target *vault.Scope
	for _, s := range scopes {
		if s.Name == name {
			copy := s
			target = &copy
			break
		}
	}
	if target == nil {
		return fmt.Errorf("scope %q not found in %s", name, scopeFile)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	secrets, err := client.FetchSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	filtered, err := vault.ApplyScope(secrets, *target)
	if err != nil {
		return fmt.Errorf("apply scope: %w", err)
	}

	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	fmt.Fprintln(os.Stdout, vault.FormatScopeReport(name, keys))

	if err := vault.WriteEnvFile(outputFile, filtered); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}
	fmt.Fprintf(os.Stdout, "wrote %d key(s) to %s\n", len(filtered), outputFile)
	return nil
}
