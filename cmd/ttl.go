package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"vaultpull/config"
	"vaultpull/vault"
)

var (
	ttlKey      string
	ttlDuration string
)

func init() {
	ttlCmd := &cobra.Command{
		Use:   "ttl",
		Short: "Manage per-key TTLs for secrets",
	}

	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Set a TTL for a secret key",
		RunE:  runTTLSet,
	}
	setCmd.Flags().StringVar(&ttlKey, "key", "", "Secret key name (required)")
	setCmd.Flags().StringVar(&ttlDuration, "ttl", "", "TTL duration e.g. 30s, 5m, 2h (required)")
	setCmd.MarkFlagRequired("key")
	setCmd.MarkFlagRequired("ttl")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List expired secret keys",
		RunE:  runTTLList,
	}

	ttlCmd.AddCommand(setCmd, listCmd)
	rootCmd.AddCommand(ttlCmd)
}

func runTTLSet(cmd *cobra.Command, args []string) error {
	_, err := config.Load()
	if err != nil {
		return err
	}
	d, err := vault.ParseTTL(ttlDuration)
	if err != nil {
		return fmt.Errorf("invalid TTL: %w", err)
	}
	m := make(vault.TTLMap)
	vault.SetTTL(m, ttlKey, d)
	fmt.Fprintf(os.Stdout, "TTL set: %s expires in %s\n", ttlKey, d)
	return nil
}

func runTTLList(cmd *cobra.Command, args []string) error {
	_, err := config.Load()
	if err != nil {
		return err
	}
	// Placeholder: in a real implementation the TTLMap would be persisted.
	m := make(vault.TTLMap)
	expired := vault.ExpiredKeys(m)
	if len(expired) == 0 {
		fmt.Println("No expired keys.")
		return nil
	}
	sort.Strings(expired)
	for _, k := range expired {
		fmt.Println(k)
	}
	return nil
}
