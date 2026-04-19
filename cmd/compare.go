package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/config"
	"vaultpull/vault"
)

func init() {
	var outputFile string
	var labelA, labelB string

	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare local .env file with secrets from Vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(outputFile, labelA, labelB)
		},
	}
	cmd.Flags().StringVar(&outputFile, "output", ".env", "Local .env file to compare against")
	cmd.Flags().StringVar(&labelA, "label-local", "local", "Label for local secrets")
	cmd.Flags().StringVar(&labelB, "label-remote", "remote", "Label for remote secrets")
	rootCmd.AddCommand(cmd)
}

func runCompare(outputFile, labelA, labelB string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	remote, err := client.FetchSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}
	local, err := vault.ParseEnvFile(outputFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("parse env file: %w", err)
	}
	result := vault.CompareSecrets(local, remote)
	fmt.Print(vault.FormatCompareResult(result, labelA, labelB))
	fmt.Printf("\nSummary: %d only-local, %d only-remote, %d differ, %d identical\n",
		len(result.OnlyInA), len(result.OnlyInB), len(result.Differ), len(result.Identical))
	return nil
}
