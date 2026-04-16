package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupEnvFile creates a timestamped backup of an existing .env file.
// Returns the backup path or an empty string if the source file does not exist.
func BackupEnvFile(envPath string) (string, error) {
	_, err := os.Stat(envPath)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", envPath, err)
	}

	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	backupName := fmt.Sprintf("%s.%s.bak", base, timestamp)
	backupPath := filepath.Join(dir, backupName)

	data, err := os.ReadFile(envPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", envPath, err)
	}

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", fmt.Errorf("write backup %s: %w", backupPath, err)
	}

	return backupPath, nil
}

// PruneBackups removes old backups for the given env file, keeping only the
// most recent `keep` backups.
func PruneBackups(envPath string, keep int) ([]string, error) {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	pattern := filepath.Join(dir, base+".*.bak")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob backups: %w", err)
	}

	if len(matches) <= keep {
		return nil, nil
	}

	// Glob returns sorted results; oldest first.
	toRemove := matches[:len(matches)-keep]
	var removed []string
	for _, f := range toRemove {
		if err := os.Remove(f); err != nil {
			return removed, fmt.Errorf("remove %s: %w", f, err)
		}
		removed = append(removed, f)
	}
	return removed, nil
}
