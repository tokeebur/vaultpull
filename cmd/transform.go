package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpull/config"
	"vaultpull/vault"
)

var transformRulesFile string

func init() {
	transformCmd := &cobra.Command{
		Use:   "transform",
		Short: "Apply value transforms to fetched secrets before writing .env",
		RunE:  runTransform,
	}
	transformCmd.Flags().StringVar(&transformRulesFile, "rules", ".vaultpull-transforms", "Path to transform rules file")
	rootCmd.AddCommand(transformCmd)
}

func runTransform(cmd *cobra.Command, args []string) error {
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
	rules, err := vault.LoadTransformFile(transformRulesFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load transform rules: %w", err)
	}
	if len(rules) > 0 {
		secrets, err = vault.ApplyTransforms(secrets, rules)
		if err != nil {
			return fmt.Errorf("apply transforms: %w", err)
		}
	}
	if err := client.WriteEnvFile(cfg.OutputFile, secrets); err != nil {
		return err
	}
	fmt.Printf("Wrote %d secrets (with transforms) to %s\n", len(secrets), cfg.OutputFile)
	return nil
}
