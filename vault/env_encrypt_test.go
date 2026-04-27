package vault

import (
	"strings"
	"testing"
)

func TestEncryptSecrets_EncryptsMatchingKeys(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret123", "APP_NAME": "myapp"}
	out, results, err := EncryptSecrets(secrets, []string{"DB_*"}, "passphrase")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out["DB_PASS"], "enc:") {
		t.Errorf("expected DB_PASS to be encrypted, got %q", out["DB_PASS"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unchanged")
	}
	if len(results) != 1 || results[0].Key != "DB_PASS" {
		t.Errorf("expected 1 result for DB_PASS")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	secrets := map[string]string{"API_KEY": "topsecret", "HOST": "localhost"}
	enc, _, err := EncryptSecrets(secrets, nil, "mypassword")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	dec, err := DecryptSecrets(enc, "mypassword")
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	for k, v := range secrets {
		if dec[k] != v {
			t.Errorf("key %q: want %q got %q", k, v, dec[k])
		}
	}
}

func TestEncryptSecrets_EmptyPassphrase(t *testing.T) {
	_, _, err := EncryptSecrets(map[string]string{"K": "v"}, nil, "")
	if err == nil {
		t.Error("expected error for empty passphrase")
	}
}

func TestEncryptSecrets_SkipsAlreadyEncrypted(t *testing.T) {
	enc, _, err := EncryptSecrets(map[string]string{"K": "val"}, nil, "pass")
	if err != nil {
		t.Fatal(err)
	}
	enc2, results, err := EncryptSecrets(enc, nil, "pass")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected no re-encryption, got %d results", len(results))
	}
	if enc2["K"] != enc["K"] {
		t.Error("value should be unchanged")
	}
}

func TestDecryptSecrets_WrongPassphrase(t *testing.T) {
	enc, _, err := EncryptSecrets(map[string]string{"K": "val"}, nil, "correct")
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecryptSecrets(enc, "wrong")
	if err == nil {
		t.Error("expected error for wrong passphrase")
	}
}

func TestFormatEncryptReport_Empty(t *testing.T) {
	out := FormatEncryptReport(nil)
	if out != "No keys encrypted." {
		t.Errorf("unexpected: %q", out)
	}
}

func TestFormatEncryptReport_WithResults(t *testing.T) {
	results := []EncryptResult{{Key: "DB_PASS", WasPlain: true}, {Key: "API_KEY", WasPlain: true}}
	out := FormatEncryptReport(results)
	if !strings.Contains(out, "Encrypted 2 key(s)") {
		t.Errorf("unexpected report: %q", out)
	}
	if !strings.Contains(out, "API_KEY") || !strings.Contains(out, "DB_PASS") {
		t.Errorf("missing key names in report")
	}
}
