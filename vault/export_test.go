package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportSecrets_Dotenv(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}

	if err := ExportSecrets(secrets, out, FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	content := string(data)
	if !strings.Contains(content, "API_KEY=") || !strings.Contains(content, "DB_PASS=") {
		t.Errorf("dotenv output missing keys: %s", content)
	}
}

func TestExportSecrets_JSON(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "secrets.json")
	secrets := map[string]string{"TOKEN": "tok_xyz"}

	if err := ExportSecrets(secrets, out, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	content := string(data)
	if !strings.Contains(content, "\"TOKEN\"") || !strings.Contains(content, "\"tok_xyz\"") {
		t.Errorf("json output malformed: %s", content)
	}
}

func TestExportSecrets_YAML(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "secrets.yaml")
	secrets := map[string]string{"HOST": "localhost"}

	if err := ExportSecrets(secrets, out, FormatYAML); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "HOST:") {
		t.Errorf("yaml output missing key: %s", string(data))
	}
}

func TestExportSecrets_UnsupportedFormat(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.txt")
	err := ExportSecrets(map[string]string{"K": "V"}, out, ExportFormat("toml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportSecrets_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	ExportSecrets(map[string]string{"X": "y"}, out, FormatDotenv)

	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
