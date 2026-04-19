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
	var groupFile string
	var group string
	var output string

	cmd := &cobra.Command{
		Use:   "group",
		Short: "Filter and export secrets by group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroup(groupFile, group, output)
		},
	}
	cmd.Flags().StringVar(&groupFile, "group-file", ".vaultgroups", "Path to group definition file")
	cmd.Flags().StringVar(&group, "group", "", "Group name to filter by (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output .env file (defaults to config output_file)")
	cmd.MarkFlagRequired("group")
	rootCmd.AddCommand(cmd)
}

func runGroup(groupFile, group, output string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	gm, err := vault.LoadGroupFile(groupFile)
	if err != nil {
		return fmt.Errorf("load group file: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.FetchSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	filtered, err := vault.FilterByGroup(secrets, gm, group)
	if err != nil {
		return fmt.Errorf("filter: %w", err)
	}

	dest := output
	if dest == "" {
		dest = cfg.OutputFile
	}

	lines := []string{}
	for k, v := range filtered {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}
	return os.WriteFile(dest, []byte(strings.Join(lines, "\n")+"\n"), 0600)
}
