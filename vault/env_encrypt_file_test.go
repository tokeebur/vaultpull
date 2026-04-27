package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncryptDecrypt_FileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	content := "DB_PASS=hunter2\nAPP_NAME=myapp\nAPI_KEY=abc123\n"
	if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	secrets, err := ParseEnvFile(envFile)
	if err != nil {
		t.Fatal(err)
	}

	enc, results, err := EncryptSecrets(secrets, []string{"DB_*", "API_*"}, "testpass")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 encrypted keys, got %d", len(results))
	}
	if !strings.HasPrefix(enc["DB_PASS"], "enc:") {
		t.Error("DB_PASS not encrypted")
	}
	if enc["APP_NAME"] != "myapp" {
		t.Error("APP_NAME should be plain")
	}

	dec, err := DecryptSecrets(enc, "testpass")
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if dec["DB_PASS"] != "hunter2" {
		t.Errorf("DB_PASS: want hunter2 got %q", dec["DB_PASS"])
	}
	if dec["API_KEY"] != "abc123" {
		t.Errorf("API_KEY: want abc123 got %q", dec["API_KEY"])
	}
}

func TestFormatDecryptReport_CountsChanges(t *testing.T) {
	orig := map[string]string{"A": "enc:abc", "B": "plain"}
	dec := map[string]string{"A": "secret", "B": "plain"}
	out := FormatDecryptReport(orig, dec)
	if !strings.Contains(out, "1") {
		t.Errorf("expected count 1 in report, got %q", out)
	}
}

func TestFormatDecryptReport_NoChanges(t *testing.T) {
	orig := map[string]string{"A": "plain"}
	dec := map[string]string{"A": "plain"}
	out := FormatDecryptReport(orig, dec)
	if out != "No keys decrypted." {
		t.Errorf("unexpected: %q", out)
	}
}
