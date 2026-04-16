package vault_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/vaultpull/vault"
)

func setupRotateServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"KEY1":"val1","KEY2":"val2"}}}`))
	}))
}

func TestRotateSecrets_NewFile(t *testing.T) {
	srv := setupRotateServer(t)
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token", "secret/data/app")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	result, err := vault.RotateSecrets(client, output)
	if err != nil {
		t.Fatalf("RotateSecrets: %v", err)
	}

	if result.OutputFile != output {
		t.Errorf("expected output %s, got %s", output, result.OutputFile)
	}
	if len(result.Diff.Added) != 2 {
		t.Errorf("expected 2 added keys, got %d", len(result.Diff.Added))
	}

	data, _ := os.ReadFile(output)
	if !strings.Contains(string(data), "KEY1=val1") {
		t.Error("output file missing KEY1=val1")
	}
}

func TestRotateSecrets_UpdatesExisting(t *testing.T) {
	srv := setupRotateServer(t)
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token", "secret/data/app")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	dir := t.TempDir()
	output := filepath.Join(dir, ".env")
	_ = os.WriteFile(output, []byte("KEY1=old\nKEY3=gone\n"), 0600)

	result, err := vault.RotateSecrets(client, output)
	if err != nil {
		t.Fatalf("RotateSecrets: %v", err)
	}

	if len(result.Diff.Changed) != 1 {
		t.Errorf("expected 1 changed key, got %d", len(result.Diff.Changed))
	}
	if len(result.Diff.Removed) != 1 {
		t.Errorf("expected 1 removed key, got %d", len(result.Diff.Removed))
	}
	if result.BackupFile == "" {
		t.Error("expected backup file path")
	}
}
