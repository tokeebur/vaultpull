package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGroupFile_NotExist(t *testing.T) {
	gm, err := LoadGroupFile("/tmp/no_such_group_file.conf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gm) != 0 {
		t.Fatalf("expected empty map")
	}
}

func TestSaveAndLoadGroupFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "groups.conf")

	original := GroupMap{
		"infra": {"DB_HOST", "DB_PORT"},
		"app":   {"APP_KEY", "APP_SECRET"},
	}
	if err := SaveGroupFile(path, original); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadGroupFile(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	for group, keys := range original {
		got := loaded[group]
		if len(got) != len(keys) {
			t.Fatalf("group %s: expected %v got %v", group, keys, got)
		}
	}
}

func TestLoadGroupFile_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "groups.conf")
	os.WriteFile(path, []byte("# comment\nweb=API_KEY,API_SECRET\n"), 0600)
	gm, err := LoadGroupFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gm["web"]; !ok {
		t.Fatal("expected group 'web'")
	}
}

func TestFilterByGroup_ReturnsMatching(t *testing.T) {
	gm := GroupMap{"db": {"DB_HOST", "DB_PORT"}}
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "OTHER": "x"}
	result, err := FilterByGroup(secrets, gm, "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["OTHER"]; ok {
		t.Fatal("OTHER should not be included")
	}
}

func TestFilterByGroup_UnknownGroup(t *testing.T) {
	gm := GroupMap{}
	_, err := FilterByGroup(map[string]string{}, gm, "missing")
	if err == nil {
		t.Fatal("expected error for unknown group")
	}
}

func TestGroupNames_Sorted(t *testing.T) {
	gm := GroupMap{"z": {}, "a": {}, "m": {}}
	names := GroupNames(gm)
	if names[0] != "a" || names[1] != "m" || names[2] != "z" {
		t.Fatalf("expected sorted names, got %v", names)
	}
}
