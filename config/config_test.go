package config

import (
	"os"
	"testing"
)

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	_, err := Load("secret/myapp", ".env")
	if err == nil {
		t.Fatal("expected error when VAULT_TOKEN is missing")
	}
}

func TestLoad_MissingSecretPath(t *testing.T) {
	os.Setenv("VAULT_TOKEN", "test-token")
	defer os.Unsetenv("VAULT_TOKEN")

	_, err := Load("", ".env")
	if err == nil {
		t.Fatal("expected error when secret path is empty")
	}
}

func TestLoad_DefaultOutputFile(t *testing.T) {
	os.Setenv("VAULT_TOKEN", "test-token")
	defer os.Unsetenv("VAULT_TOKEN")

	cfg, err := Load("secret/myapp", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected .env, got %s", cfg.OutputFile)
	}
}

func TestLoad_DefaultVaultAddr(t *testing.T) {
	os.Setenv("VAULT_TOKEN", "test-token")
	os.Unsetenv("VAULT_ADDR")
	defer os.Unsetenv("VAULT_TOKEN")

	cfg, err := Load("secret/myapp", ".env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected vault addr: %s", cfg.VaultAddr)
	}
}
