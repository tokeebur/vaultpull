package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"vaultpull/config"
	"vaultpull/vault"
)

func init() {
	var policyFile string
	var outputFile string

	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Filter synced secrets through an allow/deny policy file",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			pol, err := vault.LoadPolicyFile(policyFile)
			if err != nil {
				return fmt.Errorf("load policy: %w", err)
			}
			filtered, denied, err := vault.FilterByPolicy(secrets, pol)
			if err != nil {
				return err
			}
			if len(denied) > 0 {
				sort.Strings(denied)
				fmt.Fprintf(os.Stderr, "policy denied keys: %v\n", denied)
			}
			out := outputFile
			if out == "" {
				out = cfg.OutputFile
			}
			return vault.WriteEnvFile(out, filtered)
		},
	}

	policyCmd.Flags().StringVarP(&policyFile, "policy", "p", ".vaultpolicy", "Path to policy file")
	policyCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output .env file (overrides config)")
	rootCmd.AddCommand(policyCmd)
}
