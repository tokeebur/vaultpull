package vault

import (
	"testing"
)

func TestCastSecrets_String(t *testing.T) {
	secrets := map[string]string{"APP_NAME": "vaultpull"}
	rules := []CastRule{{Pattern: "APP_NAME", Type: CastString}}
	results, err := CastSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Cast != "vaultpull" {
		t.Errorf("expected 'vaultpull', got %v", results)
	}
}

func TestCastSecrets_Int(t *testing.T) {
	secrets := map[string]string{"PORT": "8080"}
	rules := []CastRule{{Pattern: "PORT", Type: CastInt}}
	results, err := CastSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Cast != 8080 {
		t.Errorf("expected 8080, got %v", results[0].Cast)
	}
}

func TestCastSecrets_Bool(t *testing.T) {
	secrets := map[string]string{"DEBUG": "true"}
	rules := []CastRule{{Pattern: "DEBUG", Type: CastBool}}
	results, err := CastSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Cast != true {
		t.Errorf("expected true, got %v", results[0].Cast)
	}
}

func TestCastSecrets_Float(t *testing.T) {
	secrets := map[string]string{"RATIO": "3.14"}
	rules := []CastRule{{Pattern: "RATIO", Type: CastFloat}}
	results, err := CastSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Cast.(float64) < 3.13 {
		t.Errorf("unexpected float value: %v", results[0].Cast)
	}
}

func TestCastSecrets_InvalidInt(t *testing.T) {
	secrets := map[string]string{"PORT": "not-a-number"}
	rules := []CastRule{{Pattern: "PORT", Type: CastInt}}
	_, err := CastSecrets(secrets, rules)
	if err == nil {
		t.Error("expected error for invalid int cast")
	}
}

func TestCastSecrets_WildcardPattern(t *testing.T) {
	secrets := map[string]string{"MAX_RETRIES": "5", "MAX_TIMEOUT": "30"}
	rules := []CastRule{{Pattern: "MAX_*", Type: CastInt}}
	results, err := CastSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}
