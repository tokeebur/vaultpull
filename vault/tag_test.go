package vault

import (
	"testing"
)

func TestTagSecrets_AssignsTags(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "key123"}
	tagMap := map[string][]string{
		"DB_PASS": {"database", "sensitive"},
		"API_KEY": {"api"},
	}
	result := TagSecrets(secrets, tagMap)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	// sorted: API_KEY, DB_PASS
	if result[0].Key != "API_KEY" || len(result[0].Tags) != 1 {
		t.Errorf("unexpected first entry: %+v", result[0])
	}
	if result[1].Key != "DB_PASS" || len(result[1].Tags) != 2 {
		t.Errorf("unexpected second entry: %+v", result[1])
	}
}

func TestTagSecrets_MissingTagEntry(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abc"}
	result := TagSecrets(secrets, map[string][]string{})
	if len(result[0].Tags) != 0 {
		t.Errorf("expected empty tags")
	}
}

func TestFilterByTag_ReturnsMatching(t *testing.T) {
	tagged := []TaggedSecret{
		{Key: "A", Value: "1", Tags: []string{"sensitive"}},
		{Key: "B", Value: "2", Tags: []string{"api"}},
		{Key: "C", Value: "3", Tags: []string{"sensitive", "api"}},
	}
	out := FilterByTag(tagged, "sensitive")
	if len(out) != 2 {
		t.Errorf("expected 2 matches, got %d", len(out))
	}
}

func TestFilterByTag_NoMatch(t *testing.T) {
	tagged := []TaggedSecret{
		{Key: "A", Value: "1", Tags: []string{"api"}},
	}
	out := FilterByTag(tagged, "database")
	if len(out) != 0 {
		t.Errorf("expected 0 matches")
	}
}

func TestFormatTagged(t *testing.T) {
	tagged := []TaggedSecret{
		{Key: "X", Value: "v", Tags: []string{"a", "b"}},
	}
	out := FormatTagged(tagged)
	if out != "X [a,b]\n" {
		t.Errorf("unexpected format: %q", out)
	}
}
