package vault_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/vault"
)

func TestApplyScope_IncludeAll(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc"}
	scope := vault.Scope{Name: "all"}
	result, err := vault.ApplyScope(secrets, scope)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestApplyScope_IncludePattern(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "abc"}
	scope := vault.Scope{Name: "db", Include: []string{"DB_*"}}
	result, err := vault.ApplyScope(secrets, scope)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["API_KEY"]; ok {
		t.Error("API_KEY should be excluded")
	}
}

func TestApplyScope_ExcludePattern(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PASS": "secret", "API_KEY": "abc"}
	scope := vault.Scope{Name: "safe", Exclude: []string{"*PASS*"}}
	result, err := vault.ApplyScope(secrets, scope)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["DB_PASS"]; ok {
		t.Error("DB_PASS should be excluded")
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestApplyScope_IncludeAndExclude(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "h", "DB_PASS": "p", "API_KEY": "k"}
	scope := vault.Scope{Name: "db-safe", Include: []string{"DB_*"}, Exclude: []string{"*PASS*"}}
	result, _ := vault.ApplyScope(secrets, scope)
	if len(result) != 1 || result["DB_HOST"] != "h" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestFormatScopeReport_Empty(t *testing.T) {
	out := vault.FormatScopeReport("test", []string{})
	if !strings.Contains(out, "no keys matched") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatScopeReport_WithKeys(t *testing.T) {
	out := vault.FormatScopeReport("prod", []string{"DB_HOST", "API_KEY"})
	if !strings.Contains(out, "2 key(s)") {
		t.Errorf("expected key count in output, got: %s", out)
	}
}

func TestSaveAndLoadScopeFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scopes.conf")
	scopes := []vault.Scope{
		{Name: "db", Include: []string{"DB_*"}, Exclude: []string{"*PASS*"}},
		{Name: "api", Include: []string{"API_*"}},
	}
	if err := vault.SaveScopeFile(path, scopes); err != nil {
		t.Fatal(err)
	}
	loaded, err := vault.LoadScopeFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(loaded))
	}
	if loaded[0].Name != "db" || len(loaded[0].Include) != 1 || loaded[0].Exclude[0] != "*PASS*" {
		t.Errorf("scope db mismatch: %+v", loaded[0])
	}
}

func TestLoadScopeFile_NotExist(t *testing.T) {
	scopes, err := vault.LoadScopeFile("/nonexistent/scopes.conf")
	if err != nil {
		t.Fatal(err)
	}
	if scopes != nil {
		t.Error("expected nil scopes for missing file")
	}
}

func TestSaveScopeFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scopes.conf")
	if err := vault.SaveScopeFile(path, []vault.Scope{{Name: "x"}}
	); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
