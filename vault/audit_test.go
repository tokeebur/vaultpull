package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendAuditLog_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	entry := AuditEntry{
		SecretPath: "secret/data/myapp",
		OutputFile: ".env",
		Added:      3,
		Modified:   1,
		Removed:    0,
		DryRun:     false,
	}

	if err := AppendAuditLog(logPath, entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Fatal("expected audit log file to be created")
	}
}

func TestAppendAndReadAuditLog(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	entries := []AuditEntry{
		{SecretPath: "secret/data/app1", OutputFile: ".env", Added: 2, DryRun: false},
		{SecretPath: "secret/data/app2", OutputFile: ".env.prod", Added: 0, Modified: 1, DryRun: true},
	}

	for _, e := range entries {
		if err := AppendAuditLog(logPath, e); err != nil {
			t.Fatalf("AppendAuditLog error: %v", err)
		}
	}

	got, err := ReadAuditLog(logPath)
	if err != nil {
		t.Fatalf("ReadAuditLog error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].SecretPath != "secret/data/app1" {
		t.Errorf("unexpected SecretPath: %s", got[0].SecretPath)
	}
	if !got[1].DryRun {
		t.Errorf("expected second entry DryRun=true")
	}
}

func TestReadAuditLog_NotExist(t *testing.T) {
	entries, err := ReadAuditLog("/nonexistent/path/audit.log")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}
