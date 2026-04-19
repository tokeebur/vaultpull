package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAliasFile_NotExist(t *testing.T) {
	am, err := LoadAliasFile("/tmp/no_such_alias_file.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(am) != 0 {
		t.Errorf("expected empty map, got %v", am)
	}
}

func TestSaveAndLoadAliasFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aliases.env")

	am := AliasMap{
		"DB_PASS": "DATABASE_PASSWORD",
		"API":     "API_KEY",
	}
	if err := SaveAliasFile(path, am); err != nil {
		t.Fatalf("SaveAliasFile: %v", err)
	}
	loaded, err := LoadAliasFile(path)
	if err != nil {
		t.Fatalf("LoadAliasFile: %v", err)
	}
	for alias, key := range am {
		if loaded[alias] != key {
			t.Errorf("alias %q: want %q got %q", alias, key, loaded[alias])
		}
	}
}

func TestLoadAliasFile_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aliases.env")
	os.WriteFile(path, []byte("# comment\nMY_ALIAS=REAL_KEY\n"), 0600)

	am, err := LoadAliasFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if am["MY_ALIAS"] != "REAL_KEY" {
		t.Errorf("expected MY_ALIAS=REAL_KEY, got %v", am)
	}
	if len(am) != 1 {
		t.Errorf("expected 1 entry, got %d", len(am))
	}
}

func TestApplyAliases_AddsAliasedKeys(t *testing.T) {
	secrets := map[string]string{
		"DATABASE_PASSWORD": "s3cr3t",
		"OTHER":             "value",
	}
	am := AliasMap{"DB_PASS": "DATABASE_PASSWORD"}
	out := ApplyAliases(secrets, am)

	if out["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %q", out["DB_PASS"])
	}
	if out["DATABASE_PASSWORD"] != "s3cr3t" {
		t.Errorf("original key should be preserved")
	}
}

func TestApplyAliases_MissingTarget(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	am := AliasMap{"GHOST": "NONEXISTENT"}
	out := ApplyAliases(secrets, am)
	if _, ok := out["GHOST"]; ok {
		t.Errorf("alias for missing key should not appear in output")
	}
}
