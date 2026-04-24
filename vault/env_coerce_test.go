package vault

import (
	"testing"
)

func TestCoerceSecrets_TrimSpace(t *testing.T) {
	rules := []CoerceRule{{Key: "API_KEY", Action: "trim"}}
	secrets := map[string]string{"API_KEY": "  hello  ", "OTHER": "  world  "}
	result, report, err := CoerceSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["API_KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", result["API_KEY"])
	}
	if result["OTHER"] != "  world  " {
		t.Errorf("OTHER should be unchanged, got %q", result["OTHER"])
	}
	if len(report) != 1 {
		t.Errorf("expected 1 report entry, got %d", len(report))
	}
}

func TestCoerceSecrets_ToUpper(t *testing.T) {
	rules := []CoerceRule{{Key: "ENV", Action: "upper"}}
	secrets := map[string]string{"ENV": "production"}
	result, _, err := CoerceSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["ENV"] != "PRODUCTION" {
		t.Errorf("expected 'PRODUCTION', got %q", result["ENV"])
	}
}

func TestCoerceSecrets_ToLower(t *testing.T) {
	rules := []CoerceRule{{Key: "REGION", Action: "lower"}}
	secrets := map[string]string{"REGION": "US-EAST-1"}
	result, _, err := CoerceSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["REGION"] != "us-east-1" {
		t.Errorf("expected 'us-east-1', got %q", result["REGION"])
	}
}

func TestCoerceSecrets_MissingKeySkipped(t *testing.T) {
	rules := []CoerceRule{{Key: "MISSING", Action: "trim"}}
	secrets := map[string]string{"OTHER": "value"}
	result, report, err := CoerceSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["MISSING"]; ok {
		t.Error("MISSING key should not be added")
	}
	if len(report) != 0 {
		t.Errorf("expected no report entries, got %d", len(report))
	}
}

func TestCoerceSecrets_UnknownAction(t *testing.T) {
	rules := []CoerceRule{{Key: "FOO", Action: "explode"}}
	secrets := map[string]string{"FOO": "bar"}
	_, _, err := CoerceSecrets(secrets, rules)
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestFormatCoerceReport_Empty(t *testing.T) {
	out := FormatCoerceReport(nil)
	if out != "no coercions applied" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestCoerceSecrets_DoesNotMutateOriginal(t *testing.T) {
	rules := []CoerceRule{{Key: "K", Action: "upper"}}
	secrets := map[string]string{"K": "lowercase"}
	_, _, _ = CoerceSecrets(secrets, rules)
	if secrets["K"] != "lowercase" {
		t.Error("original map should not be mutated")
	}
}
