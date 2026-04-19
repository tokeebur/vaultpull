package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendHistory_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	secrets := map[string]string{"KEY": "value"}
	if err := AppendHistory(path, "vault:secret/app", secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestAppendAndLoadHistory_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	if err := AppendHistory(path, "src1", map[string]string{"A": "1"}); err != nil {
		t.Fatal(err)
	}
	if err := AppendHistory(path, "src2", map[string]string{"A": "2", "B": "3"}); err != nil {
		t.Fatal(err)
	}
	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Secrets["A"] != "1" {
		t.Errorf("expected A=1, got %s", entries[0].Secrets["A"])
	}
	if entries[1].Source != "src2" {
		t.Errorf("expected src2, got %s", entries[1].Source)
	}
}

func TestLoadHistory_NotExist(t *testing.T) {
	_, err := LoadHistory("/nonexistent/history.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDiffHistory_DetectsChanges(t *testing.T) {
	a := HistoryEntry{Secrets: map[string]string{"X": "old", "Y": "same"}}
	b := HistoryEntry{Secrets: map[string]string{"X": "new", "Y": "same", "Z": "added"}}
	changes := DiffHistory(a, b)
	if changes["X"] != ([2]string{"old", "new"}) {
		t.Errorf("expected X change, got %v", changes["X"])
	}
	if _, ok := changes["Y"]; ok {
		t.Error("Y should not appear in diff")
	}
	if changes["Z"] != ([2]string{"", "added"}) {
		t.Errorf("expected Z added, got %v", changes["Z"])
	}
}
