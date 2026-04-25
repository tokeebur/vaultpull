package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/config"
	"vaultpull/vault"
)

func init() {
	var separator string
	var uppercase bool
	var stripPrefix string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "flatten",
		Short: "Flatten nested secret keys into a single-level env map",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFlatten(separator, stripPrefix, uppercase, dryRun)
		},
	}

	cmd.Flags().StringVar(&separator, "separator", "_", "output key separator")
	cmd.Flags().BoolVar(&uppercase, "uppercase", false, "convert keys to uppercase")
	cmd.Flags().StringVar(&stripPrefix, "strip-prefix", "", "strip prefix from keys before flattening")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print result without writing")

	rootCmd.AddCommand(cmd)
}

func runFlatten(separator, stripPrefix string, uppercase, dryRun bool) error {
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

	opts := vault.FlattenOptions{
		Separator:   separator,
		Uppercase:   uppercase,
		StripPrefix: stripPrefix,
	}

	flattened, err := vault.FlattenSecrets(secrets, opts)
	if err != nil {
		return fmt.Errorf("flatten: %w", err)
	}

	fmt.Print(vault.FormatFlattenReport(secrets, flattened))

	if dryRun {
		return nil
	}

	if err := vault.WriteEnvFile(cfg.OutputFile, flattened); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}
	fmt.Fprintf(os.Stderr, "wrote %d keys to %s\n", len(flattened), cfg.OutputFile)
	return nil
}
