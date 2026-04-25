package vault

import (
	"strings"
	"testing"
)

func TestNormalizeKeys_AlreadyUpper(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, report := NormalizeKeys(secrets, false)
	if len(report.Renamed) != 0 || len(report.Conflict) != 0 {
		t.Fatalf("expected empty report, got %+v", report)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestNormalizeKeys_RenamesLowercase(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost", "PORT": "5432"}
	out, report := NormalizeKeys(secrets, false)
	if len(report.Renamed) != 1 || report.Renamed[0] != "db_host" {
		t.Fatalf("expected db_host renamed, got %+v", report)
	}
	if out["DB_HOST"] != "localhost" {
		t.Fatalf("expected DB_HOST=localhost, got %v", out)
	}
	if _, ok := out["db_host"]; ok {
		t.Fatal("old lowercase key should not be present")
	}
}

func TestNormalizeKeys_ConflictNoOverwrite(t *testing.T) {
	// Both "api_key" and "API_KEY" are present — collision without overwrite.
	secrets := map[string]string{"api_key": "lower", "API_KEY": "upper"}
	out, report := NormalizeKeys(secrets, false)
	if len(report.Conflict) != 1 {
		t.Fatalf("expected 1 conflict, got %+v", report)
	}
	// Original uppercase value must be preserved.
	if out["API_KEY"] != "upper" {
		t.Fatalf("expected API_KEY=upper, got %s", out["API_KEY"])
	}
}

func TestNormalizeKeys_ConflictOverwrite(t *testing.T) {
	secrets := map[string]string{"api_key": "lower", "API_KEY": "upper"}
	out, report := NormalizeKeys(secrets, true)
	// With overwrite the lowercase value replaces the uppercase one.
	if len(report.Conflict) != 0 {
		t.Fatalf("expected no conflicts with overwrite, got %+v", report)
	}
	if out["API_KEY"] != "lower" {
		t.Fatalf("expected API_KEY=lower after overwrite, got %s", out["API_KEY"])
	}
}

func TestFormatUppercaseReport_Empty(t *testing.T) {
	r := UppercaseReport{}
	out := FormatUppercaseReport(r)
	if !strings.Contains(out, "no changes") {
		t.Fatalf("expected 'no changes' message, got: %s", out)
	}
}

func TestFormatUppercaseReport_WithRenamesAndConflicts(t *testing.T) {
	r := UppercaseReport{
		Renamed:  []string{"db_host"},
		Conflict: []string{"api_key"},
	}
	out := FormatUppercaseReport(r)
	if !strings.Contains(out, "renamed") {
		t.Fatalf("expected 'renamed' section, got: %s", out)
	}
	if !strings.Contains(out, "conflicts") {
		t.Fatalf("expected 'conflicts' section, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Fatalf("expected DB_HOST in output, got: %s", out)
	}
}
