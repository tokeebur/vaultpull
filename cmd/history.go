package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpull/vault"
)

func init() {
	var historyFile string
	var showDiff bool

	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show secret pull history",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := vault.LoadHistory(historyFile)
			if os.IsNotExist(err) {
				fmt.Println("No history found.")
				return nil
			}
			if err != nil {
				return fmt.Errorf("load history: %w", err)
			}
			for i, e := range entries {
				fmt.Printf("[%d] %s  source=%s  keys=%d\n", i, e.Timestamp.Format("2006-01-02 15:04:05"), e.Source, len(e.Secrets))
			}
			if showDiff && len(entries) >= 2 {
				last := entries[len(entries)-1]
				prev := entries[len(entries)-2]
				changes := vault.DiffHistory(prev, last)
				if len(changes) == 0 {
					fmt.Println("\nNo changes since previous pull.")
				} else {
					fmt.Println("\nChanges since previous pull:")
					for k, v := range changes {
						fmt.Printf("  %s: %q -> %q\n", k, v[0], v[1])
					}
				}
			}
			return nil
		},
	}

	historyCmd.Flags().StringVar(&historyFile, "file", ".vaultpull_history.json", "History file path")
	historyCmd.Flags().BoolVar(&showDiff, "diff", false, "Show diff between last two entries")
	rootCmd.AddCommand(historyCmd)
}
