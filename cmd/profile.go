package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpull/config"
	"vaultpull/vault"
)

func init() {
	var profileName string
	var profileFile string

	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Apply a named profile of secret overrides to the output env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProfile(profileName, profileFile)
		},
	}
	profileCmd.Flags().StringVarP(&profileName, "name", "n", "", "Profile name to apply (required)")
	profileCmd.Flags().StringVarP(&profileFile, "file", "f", "profiles.ini", "Path to profile file")
	profileCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(profileCmd)
}

func runProfile(name, profileFile string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	profiles, err := vault.LoadProfileFile(profileFile)
	if err != nil {
		return fmt.Errorf("load profiles: %w", err)
	}

	p, ok := profiles[name]
	if !ok {
		return fmt.Errorf("profile %q not found in %s", name, profileFile)
	}

	existing, err := vault.ParseEnvFile(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}

	merged := vault.ApplyProfile(existing, p)

	if err := vault.WriteEnvFile(cfg.OutputFile, merged); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Applied profile %q to %s (%d keys)\n", name, cfg.OutputFile, len(merged))
	return nil
}
