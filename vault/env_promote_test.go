package vault

import (
	"strings"
	"testing"
)

func TestPromoteSecrets_NewKeys(t *testing.T) {
	src := map[string]string{"DB_HOST": "prod-db", "API_KEY": "secret"}
	dst := map[string]string{}
	out, res, err := PromoteSecrets(src, dst, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if out["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST=prod-db, got %s", out["DB_HOST"])
	}
}

func TestPromoteSecrets_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"DB_HOST": "prod-db"}
	dst := map[string]string{"DB_HOST": "staging-db"}
	out, res, err := PromoteSecrets(src, dst, PromoteOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if out["DB_HOST"] != "staging-db" {
		t.Errorf("expected original value preserved")
	}
}

func TestPromoteSecrets_OverwritesExisting(t *testing.T) {
	src := map[string]string{"DB_HOST": "prod-db"}
	dst := map[string]string{"DB_HOST": "staging-db"}
	out, res, err := PromoteSecrets(src, dst, PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwrote) != 1 {
		t.Errorf("expected 1 overwrote, got %d", len(res.Overwrote))
	}
	if out["DB_HOST"] != "prod-db" {
		t.Errorf("expected overwritten value")
	}
}

func TestPromoteSecrets_DryRunDoesNotMutate(t *testing.T) {
	src := map[string]string{"NEW_KEY": "val"}
	dst := map[string]string{}
	out, res, err := PromoteSecrets(src, dst, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 {
		t.Errorf("expected 1 promoted in dry run result")
	}
	if _, ok := out["NEW_KEY"]; ok {
		t.Errorf("dry run should not write key to output")
	}
}

func TestPromoteSecrets_MissingSourceKey(t *testing.T) {
	src := map[string]string{}
	dst := map[string]string{}
	_, _, err := PromoteSecrets(src, dst, PromoteOptions{Keys: []string{"MISSING"}})
	if err == nil {
		t.Error("expected error for missing source key")
	}
}

func TestFormatPromoteReport_Empty(t *testing.T) {
	out := FormatPromoteReport(PromoteResult{})
	if !strings.Contains(out, "nothing") {
		t.Errorf("expected 'nothing' message, got: %s", out)
	}
}

func TestFormatPromoteReport_WithChanges(t *testing.T) {
	r := PromoteResult{
		Promoted:  []string{"A"},
		Overwrote: []string{"B"},
		Skipped:   []string{"C"},
	}
	out := FormatPromoteReport(r)
	for _, want := range []string{"promoted", "overwrote", "skipped"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in report", want)
		}
	}
}
