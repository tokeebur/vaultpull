package vault

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestLoadPinFile_NotExist(t *testing.T) {
	pf, err := LoadPinFile("/tmp/no_such_pin_file.pins")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pf) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestSaveAndLoadPinFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins")
	pf := PinFile{
		"DB_PASS": {Key: "DB_PASS", Value: "secret123", Comment: "pinned for prod"},
		"API_KEY": {Key: "API_KEY", Value: "abc", Comment: ""},
	}
	if err := SavePinFile(path, pf); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadPinFile(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded) != 2 {
		t.Errorf("expected 2 entries, got %d", len(loaded))
	}
	if loaded["DB_PASS"].Value != "secret123" {
		t.Errorf("wrong value for DB_PASS")
	}
	if loaded["DB_PASS"].Comment != "pinned for prod" {
		t.Errorf("comment not preserved")
	}
}

func TestApplyPins_OverridesSecrets(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "old", "HOST": "localhost"}
	pf := PinFile{"DB_PASS": {Key: "DB_PASS", Value: "pinned"}}
	out := ApplyPins(secrets, pf)
	if out["DB_PASS"] != "pinned" {
		t.Errorf("expected pinned value")
	}
	if out["HOST"] != "localhost" {
		t.Errorf("non-pinned key should be unchanged")
	}
}

func TestApplyPins_DoesNotMutateOriginal(t *testing.T) {
	secrets := map[string]string{"KEY": "original"}
	pf := PinFile{"KEY": {Key: "KEY", Value: "pinned"}}
	ApplyPins(secrets, pf)
	if secrets["KEY"] != "original" {
		t.Errorf("original map mutated")
	}
}

func TestPinnedKeys_ReturnsSorted(t *testing.T) {
	pf := PinFile{
		"Z_KEY": {Key: "Z_KEY"},
		"A_KEY": {Key: "A_KEY"},
	}
	keys := PinnedKeys(pf)
	sort.Strings(keys)
	if keys[0] != "A_KEY" || keys[1] != "Z_KEY" {
		t.Errorf("unexpected keys order: %v", keys)
	}
}

func TestSavePinFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins")
	if err := SavePinFile(path, PinFile{"K": {Key: "K", Value: "v"}}); err != nil {
		t.Fatalf("save: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
