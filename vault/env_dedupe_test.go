package vault

import (
	"strings"
	"testing"
)

func TestDedupeSecrets_NoDuplicates(t *testing.T) {
	secrets := map[string]string{"A": "foo", "B": "bar", "C": "baz"}
	out, report, err := DedupeSecrets(secrets, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report) != 0 {
		t.Errorf("expected no report entries, got %d", len(report))
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}

func TestDedupeSecrets_AlphaStrategy(t *testing.T) {
	secrets := map[string]string{"Z_KEY": "same", "A_KEY": "same", "M_KEY": "other"}
	out, report, err := DedupeSecrets(secrets, "alpha")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["A_KEY"]; !ok {
		t.Errorf("expected A_KEY to be kept")
	}
	if _, ok := out["Z_KEY"]; ok {
		t.Errorf("expected Z_KEY to be dropped")
	}
	if len(report) != 1 {
		t.Fatalf("expected 1 report entry, got %d", len(report))
	}
	if report[0].Key != "A_KEY" {
		t.Errorf("expected kept key A_KEY, got %s", report[0].Key)
	}
}

func TestDedupeSecrets_FirstStrategy(t *testing.T) {
	// "first" keeps shortest key name
	secrets := map[string]string{"LONGKEY": "dup", "K": "dup"}
	out, _, err := DedupeSecrets(secrets, "first")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["K"]; !ok {
		t.Errorf("expected K (shortest) to be kept")
	}
}

func TestDedupeSecrets_LastStrategy(t *testing.T) {
	// "last" keeps longest key name
	secrets := map[string]string{"SHORT": "dup", "MUCHLONGERKEY": "dup"}
	out, _, err := DedupeSecrets(secrets, "last")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["MUCHLONGERKEY"]; !ok {
		t.Errorf("expected MUCHLONGERKEY (longest) to be kept")
	}
}

func TestDedupeSecrets_UnknownStrategy(t *testing.T) {
	secrets := map[string]string{"A": "v"}
	_, _, err := DedupeSecrets(secrets, "random")
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestDedupeSecrets_DefaultStrategyIsAlpha(t *testing.T) {
	secrets := map[string]string{"Z": "same", "A": "same"}
	out, _, err := DedupeSecrets(secrets, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["A"]; !ok {
		t.Errorf("expected A to be kept with default alpha strategy")
	}
}

func TestFormatDedupeReport_Empty(t *testing.T) {
	msg := FormatDedupeReport(nil)
	if msg != "no duplicate values found" {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestFormatDedupeReport_ShowsDropped(t *testing.T) {
	results := []DedupeResult{
		{Key: "A_KEY", Kept: "val", Dropped: []string{"Z_KEY"}, Strategy: "alpha"},
	}
	out := FormatDedupeReport(results)
	if !strings.Contains(out, "A_KEY") || !strings.Contains(out, "Z_KEY") {
		t.Errorf("report missing expected keys: %s", out)
	}
}
