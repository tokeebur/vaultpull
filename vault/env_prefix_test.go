package vault

import (
	"strings"
	"testing"
)

func TestApplyPrefix_Add(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, results, err := ApplyPrefix(secrets, "APP_", PrefixActionAdd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_DB_HOST"] != "localhost" {
		t.Errorf("expected APP_DB_HOST=localhost, got %q", out["APP_DB_HOST"])
	}
	if out["APP_DB_PORT"] != "5432" {
		t.Errorf("expected APP_DB_PORT=5432, got %q", out["APP_DB_PORT"])
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestApplyPrefix_Remove(t *testing.T) {
	secrets := map[string]string{"APP_HOST": "localhost", "OTHER": "val"}
	out, results, err := ApplyPrefix(secrets, "APP_", PrefixActionRemove)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["OTHER"] != "val" {
		t.Errorf("expected OTHER to be preserved")
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result (only matching key), got %d", len(results))
	}
}

func TestApplyPrefix_EmptyPrefix(t *testing.T) {
	_, _, err := ApplyPrefix(map[string]string{"K": "v"}, "", PrefixActionAdd)
	if err == nil {
		t.Error("expected error for empty prefix")
	}
}

func TestApplyPrefix_UnknownAction(t *testing.T) {
	_, _, err := ApplyPrefix(map[string]string{"K": "v"}, "X_", "rewrite")
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestApplyPrefix_EmptyMap(t *testing.T) {
	out, results, err := ApplyPrefix(map[string]string{}, "X_", PrefixActionAdd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output map")
	}
	if len(results) != 0 {
		t.Errorf("expected no results")
	}
}

func TestFormatPrefixReport_Empty(t *testing.T) {
	report := FormatPrefixReport(nil)
	if report != "no keys changed" {
		t.Errorf("unexpected report: %q", report)
	}
}

func TestFormatPrefixReport_ShowsChanges(t *testing.T) {
	results := []PrefixResult{
		{OldKey: "HOST", NewKey: "APP_HOST", Action: PrefixActionAdd},
	}
	report := FormatPrefixReport(results)
	if !strings.Contains(report, "APP_HOST") {
		t.Errorf("expected report to contain APP_HOST, got: %q", report)
	}
	if !strings.Contains(report, "add") {
		t.Errorf("expected report to contain action 'add', got: %q", report)
	}
}
