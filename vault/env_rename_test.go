package vault

import (
	"testing"
)

func TestRenameSecrets_Basic(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "value1", "OTHER": "value2"}
	rules := []RenameRule{{From: "OLD_KEY", To: "NEW_KEY"}}
	out, err := RenameSecrets(secrets, rules, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if out["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", out["NEW_KEY"])
	}
	if out["OTHER"] != "value2" {
		t.Error("OTHER should be preserved")
	}
}

func TestRenameSecrets_MissingKeyNonStrict(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	rules := []RenameRule{{From: "MISSING", To: "B"}}
	out, err := RenameSecrets(secrets, rules, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["B"]; ok {
		t.Error("B should not exist")
	}
}

func TestRenameSecrets_MissingKeyStrict(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	rules := []RenameRule{{From: "MISSING", To: "B"}}
	_, err := RenameSecrets(secrets, rules, true)
	if err == nil {
		t.Fatal("expected error for missing key in strict mode")
	}
}

func TestParseRenameRules_Valid(t *testing.T) {
	lines := []string{"# comment", "", "OLD=NEW", "FOO=BAR"}
	rules, err := ParseRenameRules(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "OLD" || rules[0].To != "NEW" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseRenameRules_Invalid(t *testing.T) {
	lines := []string{"BADLINE"}
	_, err := ParseRenameRules(lines)
	if err == nil {
		t.Fatal("expected error for invalid rule")
	}
}
