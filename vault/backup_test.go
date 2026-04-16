package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBackupEnvFile_NotExist(t *testing.T) {
	path, err := BackupEnvFile("/tmp/vaultpull_nonexistent_12345.env")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if path != "" {
		t.Fatalf("expected empty path, got %s", path)
	}
}

func TestBackupEnvFile_CreatesBackup(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	original := "KEY=value\nFOO=bar\n"
	if err := os.WriteFile(envPath, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}

	backupPath, err := BackupEnvFile(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backupPath == "" {
		t.Fatal("expected a backup path")
	}
	if !strings.HasSuffix(backupPath, ".bak") {
		t.Errorf("backup path should end with .bak, got %s", backupPath)
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("could not read backup: %v", err)
	}
	if string(data) != original {
		t.Errorf("backup content mismatch: got %q", string(data))
	}
}

func TestPruneBackups_RemovesOldest(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	// Create 5 fake backup files with sortable names.
	for i := 1; i <= 5; i++ {
		name := filepath.Join(dir, ".env.2024010"+string(rune('0'+i))+"T000000Z.bak")
		if err := os.WriteFile(name, []byte("x"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	removed, err := PruneBackups(envPath, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(removed))
	}

	matches, _ := filepath.Glob(filepath.Join(dir, ".env.*.bak"))
	if len(matches) != 3 {
		t.Errorf("expected 3 remaining backups, got %d", len(matches))
	}
}
