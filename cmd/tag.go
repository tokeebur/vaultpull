package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"vaultpull/config"
	"vaultpull/vault"
)

var tagFilter string

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "List secrets filtered by tag",
		RunE:  runTag,
	}
	tagCmd.Flags().StringVar(&tagFilter, "filter", "", "Tag to filter secrets by")
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddr, cfg.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := vault.FetchSecrets(client, cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	// Build a simple tag map from KEY=tag1,tag2 env vars prefixed VAULTPULL_TAG_
	tagMap := map[string][]string{}
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "VAULTPULL_TAG_") {
			continue
		}
		parts := strings.SplitN(strings.TrimPrefix(env, "VAULTPULL_TAG_"), "=", 2)
		if len(parts) == 2 {
			tagMap[parts[0]] = strings.Split(parts[1], ",")
		}
	}

	tagged := vault.TagSecrets(secrets, tagMap)
	if tagFilter != "" {
		tagged = vault.FilterByTag(tagged, tagFilter)
	}
	fmt.Print(vault.FormatTagged(tagged))
	return nil
}
