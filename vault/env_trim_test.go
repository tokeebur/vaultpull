package vault

import (
	"strings"
	"testing"
)

func TestTrimSecrets_Space(t *testing.T) {
	secrets := map[string]string{
		"API_KEY": "  abc123  ",
		"DB_PASS": "notrailing",
	}
	rules := []TrimRule{{Key: "*", Mode: "space"}}
	out, report, err := TrimSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "abc123" {
		t.Errorf("expected trimmed value, got %q", out["API_KEY"])
	}
	if out["DB_PASS"] != "notrailing" {
		t.Errorf("expected unchanged value, got %q", out["DB_PASS"])
	}
	if len(report) != 1 || report[0].Key != "API_KEY" {
		t.Errorf("expected 1 report entry for API_KEY, got %+v", report)
	}
}

func TestTrimSecrets_Prefix(t *testing.T) {
	secrets := map[string]string{"URL": "https://example.com"}
	rules := []TrimRule{{Key: "URL", Mode: "prefix", Cutset: "https://"}}
	out, report, err := TrimSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "example.com" {
		t.Errorf("expected stripped URL, got %q", out["URL"])
	}
	if len(report) != 1 {
		t.Errorf("expected 1 report entry, got %d", len(report))
	}
}

func TestTrimSecrets_Suffix(t *testing.T) {
	secrets := map[string]string{"TOKEN": "bearer_token_"}
	rules := []TrimRule{{Key: "TOKEN", Mode: "suffix", Cutset: "_"}}
	out, _, err := TrimSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "bearer_token" {
		t.Errorf("expected suffix trimmed, got %q", out["TOKEN"])
	}
}

func TestTrimSecrets_UnknownMode(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	rules := []TrimRule{{Key: "KEY", Mode: "invalid"}}
	_, _, err := TrimSecrets(secrets, rules)
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestTrimSecrets_DoesNotMutateOriginal(t *testing.T) {
	secrets := map[string]string{"KEY": "  val  "}
	rules := []TrimRule{{Key: "KEY", Mode: "space"}}
	_, _, _ = TrimSecrets(secrets, rules)
	if secrets["KEY"] != "  val  " {
		t.Error("original secrets map was mutated")
	}
}

func TestFormatTrimReport_Empty(t *testing.T) {
	out := FormatTrimReport(nil)
	if !strings.Contains(out, "no values trimmed") {
		t.Errorf("expected empty message, got %q", out)
	}
}

func TestFormatTrimReport_ShowsChanges(t *testing.T) {
	report := []TrimReport{
		{Key: "FOO", Before: "  bar  ", After: "bar"},
	}
	out := FormatTrimReport(report)
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in report, got %q", out)
	}
	if !strings.Contains(out, "bar") {
		t.Errorf("expected value in report, got %q", out)
	}
}
