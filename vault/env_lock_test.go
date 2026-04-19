package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLockFile_NotExist(t *testing.T) {
	lf, err := LoadLockFile("/tmp/nonexistent_lock_xyz.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lf) != 0 {
		t.Errorf("expected empty map, got %d entries", len(lf))
	}
}

func TestSaveAndLoadLockFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "locks.json")

	lf := LockFile{}
	lf = LockKey(lf, "DB_PASSWORD", "alice", "rotating")
	lf = LockKey(lf, "API_KEY", "bob", "")

	if err := SaveLockFile(path, lf); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := LoadLockFile(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded) != 2 {
		t.Errorf("expected 2 entries, got %d", len(loaded))
	}
	if loaded["DB_PASSWORD"].LockedBy != "alice" {
		t.Errorf("expected alice, got %s", loaded["DB_PASSWORD"].LockedBy)
	}
}

func TestUnlockKey(t *testing.T) {
	lf := LockFile{}
	lf = LockKey(lf, "SECRET", "alice", "")
	lf, ok := UnlockKey(lf, "SECRET")
	if !ok {
		t.Error("expected ok=true")
	}
	if IsLocked(lf, "SECRET") {
		t.Error("key should be unlocked")
	}
	_, ok2 := UnlockKey(lf, "MISSING")
	if ok2 {
		t.Error("expected ok=false for missing key")
	}
}

func TestFilterLocked_RemovesLockedKeys(t *testing.T) {
	lf := LockFile{}
	lf = LockKey(lf, "LOCKED_KEY", "alice", "")

	secrets := map[string]string{
		"LOCKED_KEY": "secret",
		"FREE_KEY":   "value",
	}
	out := FilterLocked(secrets, lf)
	if _, found := out["LOCKED_KEY"]; found {
		t.Error("locked key should be filtered out")
	}
	if out["FREE_KEY"] != "value" {
		t.Error("free key should be present")
	}
}

func TestSaveLockFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "locks.json")
	if err := SaveLockFile(path, LockFile{}); err != nil {
		t.Fatalf("save: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}
