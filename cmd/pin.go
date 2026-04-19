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
	var pinFile string

	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin secret keys to fixed values, overriding Vault",
	}

	addCmd := &cobra.Command{
		Use:   "add KEY=VALUE",
		Short: "Pin a key to a value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if pinFile == "" {
				pinFile = cfg.OutputFile + ".pins"
			}
			parts := strings.SplitN(args[0], "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("expected KEY=VALUE, got %q", args[0])
			}
			pf, err := vault.LoadPinFile(pinFile)
			if err != nil {
				return err
			}
			pf[parts[0]] = vault.PinEntry{Key: parts[0], Value: parts[1]}
			return vault.SavePinFile(pinFile, pf)
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove KEY",
		Short: "Remove a pinned key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if pinFile == "" {
				pinFile = cfg.OutputFile + ".pins"
			}
			pf, err := vault.LoadPinFile(pinFile)
			if err != nil {
				return err
			}
			delete(pf, args[0])
			return vault.SavePinFile(pinFile, pf)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all pinned keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if pinFile == "" {
				pinFile = cfg.OutputFile + ".pins"
			}
			pf, err := vault.LoadPinFile(pinFile)
			if err != nil {
				return err
			}
			if len(pf) == 0 {
				fmt.Fprintln(os.Stdout, "No pinned keys.")
				return nil
			}
			for _, e := range pf {
				line := fmt.Sprintf("%s=%s", e.Key, e.Value)
				if e.Comment != "" {
					line += "  # " + e.Comment
				}
				fmt.Fprintln(os.Stdout, line)
			}
			return nil
		},
	}

	pinCmd.PersistentFlags().StringVar(&pinFile, "pin-file", "", "Path to pin file (default: <output>.pins)")
	pinCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(pinCmd)
}
