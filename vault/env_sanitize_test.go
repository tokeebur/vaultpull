package vault

import (
	"testing"
)

func TestSanitizeSecrets_TrimSpace(t *testing.T) {
	secrets := map[string]string{
		"API_KEY": "  abc123  ",
		"DB_PASS": "secret",
	}
	rules := []SanitizeRule{{Pattern: "*", Action: "trim"}}
	out, changed := SanitizeSecrets(secrets, rules)

	if out["API_KEY"] != "abc123" {
		t.Errorf("expected trimmed value, got %q", out["API_KEY"])
	}
	if out["DB_PASS"] != "secret" {
		t.Errorf("expected unchanged value, got %q", out["DB_PASS"])
	}
	if len(changed) != 1 || changed[0] != "API_KEY" {
		t.Errorf("expected [API_KEY] in changed, got %v", changed)
	}
}

func TestSanitizeSecrets_StripNonPrintable(t *testing.T) {
	secrets := map[string]string{
		"TOKEN": "abc\x01\x02def",
	}
	rules := []SanitizeRule{{Pattern: "*", Action: "strip_nonprint"}}
	out, changed := SanitizeSecrets(secrets, rules)

	if out["TOKEN"] != "abcdef" {
		t.Errorf("expected stripped value, got %q", out["TOKEN"])
	}
	if len(changed) != 1 {
		t.Errorf("expected 1 changed key, got %d", len(changed))
	}
}

func TestSanitizeSecrets_CollapseSpace(t *testing.T) {
	secrets := map[string]string{
		"DESCRIPTION": "hello   world",
	}
	rules := []SanitizeRule{{Pattern: "DESCRIPTION", Action: "collapse_space"}}
	out, changed := SanitizeSecrets(secrets, rules)

	if out["DESCRIPTION"] != "hello world" {
		t.Errorf("expected collapsed spaces, got %q", out["DESCRIPTION"])
	}
	if len(changed) != 1 {
		t.Errorf("expected 1 changed key")
	}
}

func TestSanitizeSecrets_RemoveAction(t *testing.T) {
	secrets := map[string]string{
		"TEMP_KEY": "some value",
		"KEEP_KEY": "keep",
	}
	rules := []SanitizeRule{{Pattern: "TEMP_*", Action: "remove"}}
	out, changed := SanitizeSecrets(secrets, rules)

	if out["TEMP_KEY"] != "" {
		t.Errorf("expected empty value after remove, got %q", out["TEMP_KEY"])
	}
	if out["KEEP_KEY"] != "keep" {
		t.Errorf("expected KEEP_KEY unchanged")
	}
	if len(changed) != 1 || changed[0] != "TEMP_KEY" {
		t.Errorf("expected [TEMP_KEY] changed, got %v", changed)
	}
}

func TestSanitizeSecrets_NoChange(t *testing.T) {
	secrets := map[string]string{
		"KEY": "clean",
	}
	rules := []SanitizeRule{{Pattern: "*", Action: "trim"}}
	_, changed := SanitizeSecrets(secrets, rules)

	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}

func TestSanitizeSecrets_DoesNotMutateOriginal(t *testing.T) {
	secrets := map[string]string{
		"KEY": "  value  ",
	}
	rules := []SanitizeRule{{Pattern: "*", Action: "trim"}}
	SanitizeSecrets(secrets, rules)

	if secrets["KEY"] != "  value  " {
		t.Errorf("original map was mutated")
	}
}

func TestDefaultSanitizeRules_NotEmpty(t *testing.T) {
	rules := DefaultSanitizeRules()
	if len(rules) == 0 {
		t.Error("expected non-empty default rules")
	}
}
