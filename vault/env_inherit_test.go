package vault

import (
	"strings"
	"testing"
)

func TestInheritSecrets_ParentFillsMissing(t *testing.T) {
	child := map[string]string{"APP_NAME": "myapp"}
	parent := map[string]string{"APP_NAME": "base", "LOG_LEVEL": "info"}
	rules := []InheritRule{
		{Key: "APP_NAME"},
		{Key: "LOG_LEVEL"},
	}
	out, report, err := InheritSecrets(child, parent, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected child value preserved, got %q", out["APP_NAME"])
	}
	if out["LOG_LEVEL"] != "info" {
		t.Errorf("expected parent value inherited, got %q", out["LOG_LEVEL"])
	}
	if len(report) != 2 {
		t.Errorf("expected 2 report entries, got %d", len(report))
	}
}

func TestInheritSecrets_OverrideWins(t *testing.T) {
	child := map[string]string{"DB_HOST": "localhost"}
	parent := map[string]string{"DB_HOST": "prod.db.example.com"}
	rules := []InheritRule{
		{Key: "DB_HOST", Override: true},
	}
	out, report, err := InheritSecrets(child, parent, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "prod.db.example.com" {
		t.Errorf("expected parent override, got %q", out["DB_HOST"])
	}
	if !report[0].Overridden {
		t.Error("expected Overridden=true in report")
	}
}

func TestInheritSecrets_AliasedParentKey(t *testing.T) {
	child := map[string]string{}
	parent := map[string]string{"BASE_SECRET": "s3cr3t"}
	rules := []InheritRule{
		{Key: "APP_SECRET", ParentKey: "BASE_SECRET"},
	}
	out, _, err := InheritSecrets(child, parent, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_SECRET"] != "s3cr3t" {
		t.Errorf("expected aliased parent key value, got %q", out["APP_SECRET"])
	}
}

func TestInheritSecrets_EmptyRuleKeyReturnsError(t *testing.T) {
	_, _, err := InheritSecrets(map[string]string{}, map[string]string{}, []InheritRule{{Key: ""}})
	if err == nil {
		t.Error("expected error for empty rule Key")
	}
}

func TestInheritSecrets_NilChildTreatedAsEmpty(t *testing.T) {
	parent := map[string]string{"TIMEOUT": "30s"}
	rules := []InheritRule{{Key: "TIMEOUT"}}
	out, _, err := InheritSecrets(nil, parent, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TIMEOUT"] != "30s" {
		t.Errorf("expected inherited value, got %q", out["TIMEOUT"])
	}
}

func TestFormatInheritReport_ShowsEntries(t *testing.T) {
	report := []InheritResult{
		{Key: "DB_HOST", Source: "parent", Overridden: true},
		{Key: "APP_NAME", Source: "child", Overridden: false},
	}
	out := FormatInheritReport(report)
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in report")
	}
	if !strings.Contains(out, "overridden") {
		t.Error("expected 'overridden' annotation in report")
	}
}

func TestFormatInheritReport_Empty(t *testing.T) {
	out := FormatInheritReport(nil)
	if !strings.Contains(out, "no inheritance") {
		t.Errorf("expected empty message, got %q", out)
	}
}
