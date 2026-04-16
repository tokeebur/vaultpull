package vault

import (
	"os"
	"testing"
)

func TestParseEnvFile_NotExist(t *testing.T) {
	result, err := ParseEnvFile("/tmp/nonexistent_vaultpull_test.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestParseEnvFile_ParsesKeyValues(t *testing.T) {
	f, err := os.CreateTemp("", "vaultpull_*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString("# comment\n")
	f.WriteString("FOO=bar\n")
	f.WriteString("BAZ=qux=extra\n")
	f.WriteString("\n")
	f.Close()

	result, err := ParseEnvFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", result["FOO"])
	}
	if result["BAZ"] != "qux=extra" {
		t.Errorf("expected BAZ=qux=extra, got %s", result["BAZ"])
	}
}

func TestComputeDiff(t *testing.T) {
	existing := map[string]string{
		"KEEP":    "same",
		"CHANGE":  "old",
		"REMOVED": "gone",
	}
	incoming := map[string]string{
		"KEEP":   "same",
		"CHANGE": "new",
		"ADDED":  "fresh",
	}

	diff := ComputeDiff(existing, incoming)

	if _, ok := diff.Unchanged["KEEP"]; !ok {
		t.Error("expected KEEP in Unchanged")
	}
	if diff.Changed["CHANGE"] != "new" {
		t.Errorf("expected CHANGE=new in Changed")
	}
	if _, ok := diff.Added["ADDED"]; !ok {
		t.Error("expected ADDED in Added")
	}
	if _, ok := diff.Removed["REMOVED"]; !ok {
		t.Error("expected REMOVED in Removed")
	}
}
