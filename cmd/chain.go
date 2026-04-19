package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpull/vault"
)

var chainFile string

var chainCmd = &cobra.Command{
	Use:   "chain",
	 Short: "Resolve secrets from an ordered chain of sources",
	RunE:  runChain,
}

func init() {
	chainCmd.Flags().StringVar(&chainFile, "chain-file", ".vaultchain", "Path to chain config file")
	chainCmd.Flags().StringVar(&outputFile, "output", ".env", "Output .env file")
	rootCmd.AddCommand(chainCmd)
}

func runChain(cmd *cobra.Command, args []string) error {
	cc, err := vault.LoadChainFile(chainFile)
	if err != nil {
		return fmt.Errorf("load chain file: %w", err)
	}

	sc := vault.NewSecretChain()
	for _, src := range cc.Sources {
		secrets, err := fetchFromSource(src.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: skipping source %q: %v\n", src.Name, err)
			continue
		}
		sc.AddSource(src.Name, secrets)
	}

	resolved, sources := sc.Resolve()
	for k := range resolved {
		fmt.Printf("  %s (from %s)\n", k, sources[k])
	}

	return vault.WriteEnvFile(outputFile, resolved)
}

func fetchFromSource(path string) (map[string]string, error) {
	// Attempt to parse as local .env file first, then treat as Vault path.
	if _, err := os.Stat(path); err == nil {
		return vault.ParseEnvFile(path)
	}
	return nil, fmt.Errorf("source path not found locally (Vault fetch not wired in chain cmd): %s", path)
}
