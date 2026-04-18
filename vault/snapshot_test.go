package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveSnapshot_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	secrets := map[string]string{"KEY": "value", "DB": "postgres"}
	out, err := SaveSnapshot(dir, "secret/myapp", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("snapshot file not created: %v", err)
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	secrets := map[string]string{"TOKEN": "abc123", "HOST": "localhost"}
	out, err := SaveSnapshot(dir, "secret/app", secrets)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	s, err := LoadSnapshot(out)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if s.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", s.Path)
	}
	if s.Secrets["TOKEN"] != "abc123" {
		t.Errorf("expected TOKEN=abc123, got %s", s.Secrets["TOKEN"])
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoadSnapshot_NotExist(t *testing.T) {
	_, err := LoadSnapshot(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveSnapshot_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	out, err := SaveSnapshot(dir, "secret/secure", map[string]string{"X": "y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}
