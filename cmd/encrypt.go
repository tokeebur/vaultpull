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
	var passphrase string
	var patterns []string
	var decrypt bool
	var outputFile string

	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt or decrypt secret values in a .env file using AES-GCM",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEncrypt(passphrase, patterns, decrypt, outputFile)
		},
	}

	encryptCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase for AES-GCM encryption (required)")
	encryptCmd.Flags().StringSliceVar(&patterns, "match", nil, "Key patterns to encrypt (default: all)")
	encryptCmd.Flags().BoolVar(&decrypt, "decrypt", false, "Decrypt instead of encrypt")
	encryptCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: overwrite input)")
	_ = encryptCmd.MarkFlagRequired("passphrase")

	rootCmd.AddCommand(encryptCmd)
}

func runEncrypt(passphrase string, patterns []string, decrypt bool, outputFile string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	inputFile := cfg.OutputFile
	if outputFile == "" {
		outputFile = inputFile
	}

	secrets, err := vault.ParseEnvFile(inputFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", inputFile, err)
	}

	var result map[string]string
	if decrypt {
		result, err = vault.DecryptSecrets(secrets, passphrase)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, vault.FormatDecryptReport(secrets, result))
	} else {
		var encResults []vault.EncryptResult
		result, encResults, err = vault.EncryptSecrets(secrets, patterns, passphrase)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, vault.FormatEncryptReport(encResults))
	}

	lines := make([]string, 0, len(result))
	for k, v := range result {
		lines = append(lines, k+"="+v)
	}
	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(outputFile, []byte(content), 0600)
}
