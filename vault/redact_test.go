package vault

import (
	"regexp"
	"testing"
)

func TestRedactSecrets_DefaultRules(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
		"APP_NAME":    "myapp",
	}
	out := RedactSecrets(secrets, nil)
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD redacted, got %s", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "***" {
		t.Errorf("expected API_KEY redacted, got %s", out["API_KEY"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %s", out["APP_NAME"])
	}
}

func TestRedactSecrets_CustomRules(t *testing.T) {
	rules := []RedactRule{
		{Pattern: regexp.MustCompile(`(?i)db_host`), Replacement: "<hidden>"},
	}
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	out := RedactSecrets(secrets, rules)
	if out["DB_HOST"] != "<hidden>" {
		t.Errorf("expected DB_HOST hidden, got %s", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT unchanged, got %s", out["DB_PORT"])
	}
}

func TestRedactSecrets_EmptyMap(t *testing.T) {
	out := RedactSecrets(map[string]string{}, nil)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}

func TestRedactString_ReplacesValues(t *testing.T) {
	secrets := map[string]string{"TOKEN": "tok-xyz"}
	result := RedactString("Authorization: Bearer tok-xyz", secrets)
	if result != "Authorization: Bearer ***" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestRedactString_SkipsEmptyValues(t *testing.T) {
	secrets := map[string]string{"EMPTY": ""}
	result := RedactString("some log line", secrets)
	if result != "some log line" {
		t.Errorf("unexpected result: %s", result)
	}
}
