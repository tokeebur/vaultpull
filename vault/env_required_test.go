package vault

import (
	"strings"
	"testing"
)

func TestCheckRequired_AllPresent(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	rules := []RequiredRule{
		{Key: "DB_HOST", NonEmpty: true},
		{Key: "DB_PORT", NonEmpty: false},
	}
	violations := CheckRequired(secrets, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestCheckRequired_MissingKey(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
	}
	rules := []RequiredRule{
		{Key: "DB_HOST"},
		{Key: "DB_PASS"},
	}
	violations := CheckRequired(secrets, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DB_PASS" || violations[0].Reason != "missing" {
		t.Errorf("unexpected violation: %+v", violations[0])
	}
}

func TestCheckRequired_EmptyValueNonEmpty(t *testing.T) {
	secrets := map[string]string{
		"API_KEY": "   ",
	}
	rules := []RequiredRule{
		{Key: "API_KEY", NonEmpty: true},
	}
	violations := CheckRequired(secrets, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Reason != "empty value" {
		t.Errorf("expected 'empty value', got %q", violations[0].Reason)
	}
}

func TestCheckRequired_EmptyValueNotStrict(t *testing.T) {
	secrets := map[string]string{
		"OPTIONAL_KEY": "",
	}
	rules := []RequiredRule{
		{Key: "OPTIONAL_KEY", NonEmpty: false},
	}
	violations := CheckRequired(secrets, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestFormatRequiredReport_NoViolations(t *testing.T) {
	out := FormatRequiredReport(nil)
	if out != "all required keys satisfied" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatRequiredReport_WithViolations(t *testing.T) {
	violations := []RequiredViolation{
		{Key: "Z_KEY", Reason: "missing"},
		{Key: "A_KEY", Reason: "empty value"},
	}
	out := FormatRequiredReport(violations)
	if !strings.Contains(out, "2 required key violation") {
		t.Errorf("expected count in output, got: %q", out)
	}
	lines := strings.Split(out, "\n")
	// sorted: A_KEY should appear before Z_KEY
	if !strings.Contains(lines[1], "A_KEY") {
		t.Errorf("expected A_KEY first, got: %q", lines[1])
	}
	if !strings.Contains(lines[2], "Z_KEY") {
		t.Errorf("expected Z_KEY second, got: %q", lines[2])
	}
}
