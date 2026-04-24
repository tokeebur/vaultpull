package vault

import (
	"testing"
)

func TestApplyPlaceholders_NoMatch(t *testing.T) {
	secrets := map[string]string{"APP_NAME": "myapp", "VERSION": "1.0"}
	rules := DefaultPlaceholderRules()
	out, err := ApplyPlaceholders(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" || out["VERSION"] != "1.0" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestApplyPlaceholders_ReplacesToken(t *testing.T) {
	secrets := map[string]string{"GITHUB_TOKEN": "ghp_abc123", "APP_NAME": "myapp"}
	rules := DefaultPlaceholderRules()
	out, err := ApplyPlaceholders(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GITHUB_TOKEN"] != "<secret>" {
		t.Errorf("expected <secret>, got %q", out["GITHUB_TOKEN"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", out["APP_NAME"])
	}
}

func TestApplyPlaceholders_ReplacesURL(t *testing.T) {
	secrets := map[string]string{"DATABASE_URL": "postgres://user:pass@host/db"}
	rules := DefaultPlaceholderRules()
	out, err := ApplyPlaceholders(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_URL"] != "<url>" {
		t.Errorf("expected <url>, got %q", out["DATABASE_URL"])
	}
}

func TestApplyPlaceholders_CustomRule(t *testing.T) {
	secrets := map[string]string{"MY_CERT": "-----BEGIN CERT-----"}
	rules := []PlaceholderRule{{Pattern: "*_CERT", Placeholder: "<cert>"}}
	out, err := ApplyPlaceholders(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MY_CERT"] != "<cert>" {
		t.Errorf("expected <cert>, got %q", out["MY_CERT"])
	}
}

func TestApplyPlaceholders_InvalidPattern(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	rules := []PlaceholderRule{{Pattern: "[", Placeholder: "<x>"}}
	_, err := ApplyPlaceholders(secrets, rules)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestListPlaceholderKeys_ReturnsSorted(t *testing.T) {
	original := map[string]string{"Z_TOKEN": "tok", "A_SECRET": "sec", "NAME": "app"}
	rules := DefaultPlaceholderRules()
	placeholdered, _ := ApplyPlaceholders(original, rules)
	keys := ListPlaceholderKeys(original, placeholdered)
	if len(keys) != 2 {
		t.Fatalf("expected 2 replaced keys, got %d: %v", len(keys), keys)
	}
	if keys[0] != "A_SECRET" || keys[1] != "Z_TOKEN" {
		t.Errorf("expected sorted [A_SECRET, Z_TOKEN], got %v", keys)
	}
}
