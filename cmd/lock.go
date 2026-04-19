package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpull/vault"
)

func init() {
	var lockFile string
	var reason string
	var user string

	lockCmd := &cobra.Command{
		Use:   "lock <key>",
		Short: "Lock a secret key to prevent overwrites",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			lf, err := vault.LoadLockFile(lockFile)
			if err != nil {
				return fmt.Errorf("load lock file: %w", err)
			}
			lf = vault.LockKey(lf, key, user, reason)
			if err := vault.SaveLockFile(lockFile, lf); err != nil {
				return fmt.Errorf("save lock file: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Locked key: %s (by %s)\n", key, user)
			return nil
		},
	}

	unlockCmd := &cobra.Command{
		Use:   "unlock <key>",
		Short: "Unlock a previously locked secret key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			lf, err := vault.LoadLockFile(lockFile)
			if err != nil {
				return fmt.Errorf("load lock file: %w", err)
			}
			lf, ok := vault.UnlockKey(lf, key)
			if !ok {
				return fmt.Errorf("key %q is not locked", key)
			}
			if err := vault.SaveLockFile(lockFile, lf); err != nil {
				return fmt.Errorf("save lock file: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Unlocked key: %s\n", key)
			return nil
		},
	}

	for _, c := range []*cobra.Command{lockCmd, unlockCmd} {
		c.Flags().StringVar(&lockFile, "lock-file", ".vaultlocks", "Path to lock file")
	}
	lockCmd.Flags().StringVar(&reason, "reason", "", "Reason for locking")
	lockCmd.Flags().StringVar(&user, "user", os.Getenv("USER"), "User locking the key")

	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(unlockCmd)
}
