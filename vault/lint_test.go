package vault

import (
	"testing"
)

func TestLintSecrets_NoViolations(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
	}
	results := LintSecrets(secrets, DefaultRules())
	if len(results) != 0 {
		t.Fatalf("expected no violations, got %d: %+v", len(results), results)
	}
}

func TestLintSecrets_EmptyValue(t *testing.T) {
	secrets := map[string]string{
		"MY_KEY": "",
	}
	results := LintSecrets(secrets, DefaultRules())
	if !containsRule(results, "no-empty-value") {
		t.Fatal("expected no-empty-value violation")
	}
}

func TestLintSecrets_LowercaseKey(t *testing.T) {
	secrets := map[string]string{
		"my_key": "value",
	}
	results := LintSecrets(secrets, DefaultRules())
	if !containsRule(results, "no-lowercase-key") {
		t.Fatal("expected no-lowercase-key violation")
	}
}

func TestLintSecrets_SpaceInKey(t *testing.T) {
	secrets := map[string]string{
		"MY KEY": "value",
	}
	results := LintSecrets(secrets, DefaultRules())
	if !containsRule(results, "no-space-in-key") {
		t.Fatal("expected no-space-in-key violation")
	}
}

func TestLintSecrets_MultipleViolations(t *testing.T) {
	secrets := map[string]string{
		"bad key": "",
	}
	results := LintSecrets(secrets, DefaultRules())
	if len(results) < 2 {
		t.Fatalf("expected at least 2 violations, got %d", len(results))
	}
}

func containsRule(results []LintResult, rule string) bool {
	for _, r := range results {
		if r.Rule == rule {
			return true
		}
	}
	return false
}
