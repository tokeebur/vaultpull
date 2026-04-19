package vault

import (
	"strings"
	"testing"
)

func TestCompareSecrets_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"FOO": "1"}
	r := CompareSecrets(a, b)
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "BAR" {
		t.Errorf("expected BAR only in A, got %v", r.OnlyInA)
	}
}

func TestCompareSecrets_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"FOO": "1", "BAZ": "3"}
	r := CompareSecrets(a, b)
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "BAZ" {
		t.Errorf("expected BAZ only in B, got %v", r.OnlyInB)
	}
}

func TestCompareSecrets_Differ(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"FOO": "2"}
	r := CompareSecrets(a, b)
	if len(r.Differ) != 1 || r.Differ[0] != "FOO" {
		t.Errorf("expected FOO to differ, got %v", r.Differ)
	}
}

func TestCompareSecrets_Identical(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"FOO": "1"}
	r := CompareSecrets(a, b)
	if len(r.Identical) != 1 || r.Identical[0] != "FOO" {
		t.Errorf("expected FOO identical, got %v", r.Identical)
	}
}

func TestCompareSecrets_Empty(t *testing.T) {
	r := CompareSecrets(map[string]string{}, map[string]string{})
	if len(r.OnlyInA)+len(r.OnlyInB)+len(r.Differ)+len(r.Identical) != 0 {
		t.Error("expected empty result")
	}
}

func TestFormatCompareResult(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "x"}
	b := map[string]string{"FOO": "2", "BAZ": "3"}
	r := CompareSecrets(a, b)
	out := FormatCompareResult(r, "local", "remote")
	if !strings.Contains(out, "< BAR") {
		t.Errorf("expected BAR only in local, got:\n%s", out)
	}
	if !strings.Contains(out, "> BAZ") {
		t.Errorf("expected BAZ only in remote, got:\n%s", out)
	}
	if !strings.Contains(out, "~ FOO") {
		t.Errorf("expected FOO to differ, got:\n%s", out)
	}
}
