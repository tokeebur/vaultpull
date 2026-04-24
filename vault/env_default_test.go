package vault

import (
	"testing"
)

func TestApplyDefaults_FillsMissing(t *testing.T) {
	secrets := map[string]string{"EXISTING": "val"}
	rules := []DefaultRule{{Key: "MISSING", Value: "default_val"}}
	out, applied := ApplyDefaults(secrets, rules)
	if out["MISSING"] != "default_val" {
		t.Errorf("expected default_val, got %q", out["MISSING"])
	}
	if len(applied) != 1 || applied[0] != "MISSING" {
		t.Errorf("expected MISSING in applied, got %v", applied)
	}
}

func TestApplyDefaults_DoesNotOverrideExisting(t *testing.T) {
	secrets := map[string]string{"KEY": "real"}
	rules := []DefaultRule{{Key: "KEY", Value: "fallback"}}
	out, applied := ApplyDefaults(secrets, rules)
	if out["KEY"] != "real" {
		t.Errorf("expected real, got %q", out["KEY"])
	}
	if len(applied) != 0 {
		t.Errorf("expected no applied defaults, got %v", applied)
	}
}

func TestApplyDefaults_OnEmpty_ReplacesEmptyValue(t *testing.T) {
	secrets := map[string]string{"KEY": ""}
	rules := []DefaultRule{{Key: "KEY", Value: "filled", OnEmpty: true}}
	out, applied := ApplyDefaults(secrets, rules)
	if out["KEY"] != "filled" {
		t.Errorf("expected filled, got %q", out["KEY"])
	}
	if len(applied) != 1 {
		t.Errorf("expected 1 applied, got %v", applied)
	}
}

func TestApplyDefaults_DoesNotMutateOriginal(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	rules := []DefaultRule{{Key: "B", Value: "2"}}
	ApplyDefaults(secrets, rules)
	if _, ok := secrets["B"]; ok {
		t.Error("original map should not be mutated")
	}
}

func TestFormatDefaultReport_Empty(t *testing.T) {
	result := FormatDefaultReport(nil)
	if result != "No defaults applied." {
		t.Errorf("unexpected: %q", result)
	}
}

func TestFormatDefaultReport_WithEntries(t *testing.T) {
	result := FormatDefaultReport([]string{"FOO", "BAR"})
	if result == "" {
		t.Error("expected non-empty report")
	}
	for _, key := range []string{"FOO", "BAR", "2 default"} {
		if !containsStr(result, key) {
			t.Errorf("expected %q in report: %s", key, result)
		}
	}
}

func TestParseDefaultRules_Valid(t *testing.T) {
	lines := []string{"FOO=bar", "BAZ?=qux", "# comment", ""}
	rules, err := ParseDefaultRules(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Key != "FOO" || rules[0].Value != "bar" || rules[0].OnEmpty {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
	if rules[1].Key != "BAZ" || rules[1].Value != "qux" || !rules[1].OnEmpty {
		t.Errorf("unexpected rule[1]: %+v", rules[1])
	}
}

func TestParseDefaultRules_Invalid(t *testing.T) {
	lines := []string{"NODIVIDER"}
	_, err := ParseDefaultRules(lines)
	if err == nil {
		t.Error("expected error for invalid rule")
	}
}
