package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"vaultpull/config"
	"vaultpull/vault"
)

var schemaFile string

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate fetched secrets against a schema file",
		RunE:  runSchema,
	}
	schemaCmd.Flags().StringVar(&schemaFile, "schema", ".vaultschema", "Path to schema definition file")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	secrets, err := client.FetchSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}
	schema, err := vault.LoadSchemaFile(schemaFile)
	if err != nil {
		return fmt.Errorf("load schema: %w", err)
	}
	violations := vault.ValidateAgainstSchema(secrets, schema)
	if len(violations) == 0 {
		fmt.Println("schema validation passed")
		return nil
	}
	fmt.Fprintln(os.Stderr, "schema violations:")
	fmt.Fprintln(os.Stderr, strings.Join(violations, "\n"))
	return fmt.Errorf("%d violation(s) found", len(violations))
}
