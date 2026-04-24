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
	var srcFile string
	var dstFile string
	var overwrite bool
	var dryRun bool
	var keys string

	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote secrets from one env file into another",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.Load()
			if err != nil {
				return err
			}

			src, err := vault.ParseEnvFile(srcFile)
			if err != nil {
				return fmt.Errorf("reading source file: %w", err)
			}

			dst, err := vault.ParseEnvFile(dstFile)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("reading destination file: %w", err)
			}
			if dst == nil {
				dst = map[string]string{}
			}

			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					k = strings.TrimSpace(k)
					if k != "" {
						keyList = append(keyList, k)
					}
				}
			}

			opts := vault.PromoteOptions{
				Overwrite: overwrite,
				DryRun:    dryRun,
				Keys:      keyList,
			}

			out, result, err := vault.PromoteSecrets(src, dst, opts)
			if err != nil {
				return err
			}

			fmt.Print(vault.FormatPromoteReport(result))

			if dryRun {
				fmt.Println("[dry-run] no changes written")
				return nil
			}

			return vault.WriteEnvFile(dstFile, out)
		},
	}

	cmd.Flags().StringVar(&srcFile, "src", "", "source .env file (required)")
	cmd.Flags().StringVar(&dstFile, "dst", "", "destination .env file (required)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in destination")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without writing")
	cmd.Flags().StringVar(&keys, "keys", "", "comma-separated list of keys to promote (default: all)")
	_ = cmd.MarkFlagRequired("src")
	_ = cmd.MarkFlagRequired("dst")

	rootCmd.AddCommand(cmd)
}
