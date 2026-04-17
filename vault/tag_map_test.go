package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTagMap_NotExist(t *testing.T) {
	m, err := LoadTagMap("/nonexistent/tags.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestSaveAndLoadTagMap(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.txt")

	input := map[string][]string{
		"DB_PASS": {"sensitive", "database"},
		"API_KEY": {"api"},
	}
	if err := SaveTagMap(path, input); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := LoadTagMap(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded))
	}
	if len(loaded["DB_PASS"]) != 2 {
		t.Errorf("expected 2 tags for DB_PASS")
	}
}

func TestLoadTagMap_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.txt")
	content := "# this is a comment\nTOKEN=auth\n"
	os.WriteFile(path, []byte(content), 0600)

	m, err := LoadTagMap(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m["TOKEN"]; !ok {
		t.Errorf("expected TOKEN in map")
	}
	if len(m) != 1 {
		t.Errorf("expected 1 entry, got %d", len(m))
	}
}
