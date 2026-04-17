package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/config"
	"github.com/vaultpull/vault"
)

var (
	templateFile string
	templateOutput string
)

func init() {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Render a template file using secrets from Vault",
		RunE:  runTemplate,
	}

	templateCmd.Flags().StringVarP(&templateFile, "template", "t", "", "Path to the template file (required)")
	templateCmd.Flags().StringVarP(&templateOutput, "output", "o", "", "Path for rendered output file (required)")
	_ = templateCmd.MarkFlagRequired("template")
	_ = templateCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	secrets, err := vault.FetchSecrets(client, cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetching secrets: %w", err)
	}

	if err := vault.RenderTemplate(templateFile, templateOutput, secrets); err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Template rendered to %s\n", templateOutput)
	return nil
}
