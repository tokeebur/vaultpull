package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateSecretKeys_AllValid(t *testing.T) {
	secrets := map[string]interface{}{
		"DB_HOST": "localhost",
		"API_KEY": "abc123",
		"PORT":    "8080",
	}
	if err := ValidateSecretKeys(secrets); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSecretKeys_InvalidKeys(t *testing.T) {
	secrets := map[string]interface{}{
		"valid_KEY": "ok",
		"1INVALID":  "bad",
		"also-bad":  "bad",
	}
	err := ValidateSecretKeys(secrets)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.InvalidKeys) != 3 {
		t.Errorf("expected 3 invalid keys, got %d: %v", len(ve.InvalidKeys), ve.InvalidKeys)
	}
}

func TestValidateEnvFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "DB_HOST=localhost\nAPI_KEY=secret\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateEnvFile(path); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateEnvFile_NotExist(t *testing.T) {
	err := ValidateEnvFile("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidateSecretKeys_Empty(t *testing.T) {
	if err := ValidateSecretKeys(map[string]interface{}{}); err != nil {
		t.Errorf("unexpected error for empty map: %v", err)
	}
}
