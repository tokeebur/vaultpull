package vault

import (
	"strings"
	"testing"
)

func TestSortSecrets_Alpha(t *testing.T) {
	secrets := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	result, err := SortSecrets(secrets, SortAlpha)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Key != "APPLE" || result[1].Key != "MANGO" || result[2].Key != "ZEBRA" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSortSecrets_AlphaDesc(t *testing.T) {
	secrets := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	result, err := SortSecrets(secrets, SortAlphaDesc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Key != "ZEBRA" || result[1].Key != "MANGO" || result[2].Key != "APPLE" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSortSecrets_ByLength(t *testing.T) {
	secrets := map[string]string{"AB": "1", "ABCDE": "2", "ABC": "3"}
	result, err := SortSecrets(secrets, SortByLength)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Key != "AB" || result[1].Key != "ABC" || result[2].Key != "ABCDE" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSortSecrets_UnknownOrder(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	_, err := SortSecrets(secrets, SortOrder("bogus"))
	if err == nil {
		t.Fatal("expected error for unknown sort order")
	}
}

func TestSortSecrets_EmptyMap(t *testing.T) {
	result, err := SortSecrets(map[string]string{}, SortAlpha)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result")
	}
}

func TestFormatSorted(t *testing.T) {
	sorted := []SortedSecret{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}
	out := FormatSorted(sorted)
	if !strings.Contains(out, "FOO=bar") || !strings.Contains(out, "BAZ=qux") {
		t.Errorf("unexpected output: %q", out)
	}
}
