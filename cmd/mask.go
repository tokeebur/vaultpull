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
	var showLast int
	var envFile string

	maskCmd := &cobra.Command{
		Use:   "mask",
		Short: "Print secrets with sensitive values masked",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if envFile == "" {
				envFile = cfg.OutputFile
			}
			secrets, err := vault.ParseEnvFile(envFile)
			if err != nil {
				return fmt.Errorf("reading env file: %w", err)
			}
			rules := vault.DefaultMaskRules()
			if showLast >= 0 {
				// allow override via flag for all rules
				for i := range rules {
					if showLast == 0 {
						rules[i].ShowLast = 0
					}
				}
			}
			masked := vault.MaskSecrets(secrets, rules)
			keys := make([]string, 0, len(masked))
			for k := range masked {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, masked[k])
			}
			return nil
		},
	}

	maskCmd.Flags().StringVar(&envFile, "file", "", "env file to mask (default from config)")
	maskCmd.Flags().IntVar(&showLast, "show-last", -1, "override show-last chars for all rules (0 = full mask)")
	rootCmd.AddCommand(maskCmd)
}
