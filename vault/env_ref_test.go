package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildRefIndex_MapsKeys(t *testing.T) {
	refs := []RefEntry{
		{Key: "DB_PASS", Source: "vault:secret/db#pass"},
		{Key: "API_KEY", Source: "env:API_KEY"},
	}
	idx := BuildRefIndex(refs)
	if idx["DB_PASS"].Source != "vault:secret/db#pass" {
		t.Errorf("unexpected source: %s", idx["DB_PASS"].Source)
	}
	if _, ok := idx["MISSING"]; ok {
		t.Error("expected missing key to be absent")
	}
}

func TestResolveRefs_LiteralSource(t *testing.T) {
	secrets := map[string]string{"FOO": "old"}
	refs := []RefEntry{{Key: "FOO", Source: "literal:new_value"}}
	out, err := ResolveRefs(secrets, refs)
	if err != nil {
		t.Fatal(err)
	}
	if out["FOO"] != "new_value" {
		t.Errorf("expected new_value, got %s", out["FOO"])
	}
}

func TestResolveRefs_NonLiteralUnchanged(t *testing.T) {
	secrets := map[string]string{"FOO": "original"}
	refs := []RefEntry{{Key: "FOO", Source: "vault:secret/app#FOO"}}
	out, _ := ResolveRefs(secrets, refs)
	if out["FOO"] != "original" {
		t.Errorf("expected original, got %s", out["FOO"])
	}
}

func TestFormatRefReport_Empty(t *testing.T) {
	out := FormatRefReport(nil)
	if !strings.Contains(out, "no refs") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatRefReport_Sorted(t *testing.T) {
	refs := []RefEntry{
		{Key: "Z_KEY", Source: "literal:z"},
		{Key: "A_KEY", Source: "literal:a", Note: "primary"},
	}
	out := FormatRefReport(refs)
	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA > idxZ {
		t.Error("expected A_KEY before Z_KEY")
	}
	if !strings.Contains(out, "primary") {
		t.Error("expected note in output")
	}
}

func TestSaveAndLoadRefFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "refs.conf")
	refs := []RefEntry{
		{Key: "DB_PASS", Source: "vault:secret/db#pass", Note: "main db"},
		{Key: "API_KEY", Source: "env:API_KEY"},
	}
	if err := SaveRefFile(path, refs); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadRefFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded))
	}
	if loaded[0].Note != "main db" {
		t.Errorf("note not preserved: %s", loaded[0].Note)
	}
}

func TestLoadRefFile_NotExist(t *testing.T) {
	entries, err := LoadRefFile("/nonexistent/path/refs.conf")
	if err != nil {
		t.Fatal(err)
	}
	if entries != nil {
		t.Error("expected nil for missing file")
	}
}

func TestSaveRefFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "refs.conf")
	_ = SaveRefFile(path, nil)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
