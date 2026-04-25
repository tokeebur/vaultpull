package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/vault"
)

func init() {
	var refFile string
	var outputFile string

	refCmd := &cobra.Command{
		Use:   "ref",
		Short: "Manage cross-references between secret keys and external sources",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all defined refs",
		RunE: func(cmd *cobra.Command, args []string) error {
			refs, err := vault.LoadRefFile(refFile)
			if err != nil {
				return fmt.Errorf("load ref file: %w", err)
			}
			fmt.Print(vault.FormatRefReport(refs))
			return nil
		},
	}

	addCmd := &cobra.Command{
		Use:   "add KEY SOURCE [note]",
		Short: "Add or update a ref entry",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			refs, err := vault.LoadRefFile(refFile)
			if err != nil {
				return fmt.Errorf("load ref file: %w", err)
			}
			entry := vault.RefEntry{Key: args[0], Source: args[1]}
			if len(args) == 3 {
				entry.Note = args[2]
			}
			// Replace existing or append
			updated := false
			for i, r := range refs {
				if r.Key == entry.Key {
					refs[i] = entry
					updated = true
					break
				}
			}
			if !updated {
				refs = append(refs, entry)
			}
			if err := vault.SaveRefFile(refFile, refs); err != nil {
				return fmt.Errorf("save ref file: %w", err)
			}
			fmt.Fprintf(os.Stdout, "ref saved: %s -> %s\n", entry.Key, entry.Source)
			return nil
		},
	}

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply literal refs to the output env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			refs, err := vault.LoadRefFile(refFile)
			if err != nil {
				return fmt.Errorf("load ref file: %w", err)
			}
			existing, err := vault.ParseEnvFile(outputFile)
			if err != nil {
				return fmt.Errorf("parse env file: %w", err)
			}
			resolved, err := vault.ResolveRefs(existing, refs)
			if err != nil {
				return err
			}
			if err := vault.WriteEnvFile(outputFile, resolved); err != nil {
				return fmt.Errorf("write env file: %w", err)
			}
			fmt.Fprintf(os.Stdout, "refs applied to %s\n", outputFile)
			return nil
		},
	}

	refCmd.PersistentFlags().StringVar(&refFile, "ref-file", ".vaultpull-refs", "path to ref definition file")
	applyCmd.Flags().StringVar(&outputFile, "output", ".env", "env file to apply refs into")

	refCmd.AddCommand(listCmd, addCmd, applyCmd)
	rootCmd.AddCommand(refCmd)
}
