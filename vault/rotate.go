package vault

import (
	"fmt"
	"os"
	"time"
)

// RotateResult holds the outcome of a rotation operation.
type RotateResult struct {
	OutputFile string
	BackupFile string
	Diff       DiffResult
	Timestamp  time.Time
}

// RotateSecrets fetches fresh secrets from Vault, diffs against the existing
// .env file, creates a backup, writes the new file, and appends an audit entry.
func RotateSecrets(client *Client, outputFile string) (*RotateResult, error) {
	secrets, err := client.FetchSecrets()
	if err != nil {
		return nil, fmt.Errorf("rotate: fetch secrets: %w", err)
	}

	existing, err := ParseEnvFile(outputFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("rotate: parse existing env: %w", err)
	}

	diff := ComputeDiff(existing, secrets)

	backupPath, err := BackupEnvFile(outputFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("rotate: backup: %w", err)
	}

	if err := client.WriteEnvFile(secrets, outputFile); err != nil {
		return nil, fmt.Errorf("rotate: write env: %w", err)
	}

	now := time.Now().UTC()
	entry := AuditEntry{
		Timestamp:  now,
		Operation:  "rotate",
		OutputFile: outputFile,
		Added:      len(diff.Added),
		Removed:    len(diff.Removed),
		Changed:    len(diff.Changed),
	}
	if aerr := AppendAuditLog(outputFile+".audit.log", entry); aerr != nil {
		fmt.Fprintf(os.Stderr, "warning: audit log: %v\n", aerr)
	}

	return &RotateResult{
		OutputFile: outputFile,
		BackupFile: backupPath,
		Diff:       diff,
		Timestamp:  now,
	}, nil
}
