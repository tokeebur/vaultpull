package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpull/config"
	"vaultpull/vault"
)

func init() {
	var annotationFile string
	var addKey string
	var addNote string

	cmd := &cobra.Command{
		Use:   "annotate",
		Short: "View or add annotations to secret keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}

			am, err := vault.LoadAnnotationFile(annotationFile)
			if err != nil {
				return fmt.Errorf("load annotations: %w", err)
			}

			if addKey != "" && addNote != "" {
				am = vault.AddAnnotation(am, addKey, addNote)
				if err := vault.SaveAnnotationFile(annotationFile, am); err != nil {
					return fmt.Errorf("save annotations: %w", err)
				}
				fmt.Printf("Annotation added: %s -> %s\n", addKey, addNote)
				return nil
			}

			envData, err := vault.ParseEnvFile(cfg.OutputFile)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("read env file: %w", err)
			}

			out := vault.FormatAnnotations(envData, am)
			if out == "" {
				fmt.Println("No annotations found.")
				return nil
			}
			fmt.Print(out)
			return nil
		},
	}

	cmd.Flags().StringVar(&annotationFile, "file", ".annotations", "Path to annotation file")
	cmd.Flags().StringVar(&addKey, "key", "", "Secret key to annotate")
	cmd.Flags().StringVar(&addNote, "note", "", "Note to attach to the key")

	rootCmd.AddCommand(cmd)
}
