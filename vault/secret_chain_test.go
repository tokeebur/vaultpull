package vault

import (
	"strings"
	"testing"
)

func TestNewSecretChain_Empty(t *testing.T) {
	sc := NewSecretChain()
	if len(sc.Entries) != 0 {
		t.Fatal("expected empty chain")
	}
}

func TestAddSource_AddsEntries(t *testing.T) {
	sc := NewSecretChain()
	sc.AddSource("vault", map[string]string{"FOO": "bar"})
	if len(sc.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(sc.Entries))
	}
	if sc.Entries[0].Source != "vault" || sc.Entries[0].Key != "FOO" {
		t.Fatal("unexpected entry content")
	}
}

func TestResolve_LastWins(t *testing.T) {
	sc := NewSecretChain()
	sc.AddSource("base", map[string]string{"FOO": "base_val", "BAR": "bar_val"})
	sc.AddSource("override", map[string]string{"FOO": "override_val"})
	resolved, sources := sc.Resolve()
	if resolved["FOO"] != "override_val" {
		t.Errorf("expected override_val, got %s", resolved["FOO"])
	}
	if sources["FOO"] != "override" {
		t.Errorf("expected source=override, got %s", sources["FOO"])
	}
	if resolved["BAR"] != "bar_val" {
		t.Errorf("expected bar_val, got %s", resolved["BAR"])
	}
}

func TestExplain_ShowsChain(t *testing.T) {
	sc := NewSecretChain()
	sc.AddSource("base", map[string]string{"FOO": "v1"})
	sc.AddSource("env", map[string]string{"FOO": "v2"})
	out := sc.Explain("FOO")
	if !strings.Contains(out, "base") || !strings.Contains(out, "env") {
		t.Errorf("expected both sources in explain output: %s", out)
	}
}

func TestExplain_MissingKey(t *testing.T) {
	sc := NewSecretChain()
	out := sc.Explain("MISSING")
	if !strings.Contains(out, "no entries found") {
		t.Errorf("expected no entries message, got: %s", out)
	}
}
