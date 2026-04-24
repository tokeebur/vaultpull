package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadFreezeFile_NotExist(t *testing.T) {
	ff, err := LoadFreezeFile("/nonexistent/freeze.conf")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(ff) != 0 {
		t.Errorf("expected empty freeze file, got %v", ff)
	}
}

func TestSaveAndLoadFreezeFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "freeze.conf")

	ff := FreezeFile{
		"DB_PASSWORD": {Key: "DB_PASSWORD", FrozenAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), FrozenBy: "alice"},
		"API_KEY":     {Key: "API_KEY", FrozenAt: time.Date(2024, 2, 20, 12, 0, 0, 0, time.UTC), FrozenBy: "bob"},
	}

	if err := SaveFreezeFile(path, ff); err != nil {
		t.Fatalf("SaveFreezeFile: %v", err)
	}

	loaded, err := LoadFreezeFile(path)
	if err != nil {
		t.Fatalf("LoadFreezeFile: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded))
	}
	if loaded["DB_PASSWORD"].FrozenBy != "alice" {
		t.Errorf("expected FrozenBy=alice, got %s", loaded["DB_PASSWORD"].FrozenBy)
	}
	if loaded["API_KEY"].FrozenBy != "bob" {
		t.Errorf("expected FrozenBy=bob, got %s", loaded["API_KEY"].FrozenBy)
	}
}

func TestLoadFreezeFile_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "freeze.conf")
	content := "# comment\nSECRET_KEY=2024-01-01T00:00:00Z|ci\n# another comment\n"
	os.WriteFile(path, []byte(content), 0600)

	ff, err := LoadFreezeFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ff) != 1 {
		t.Errorf("expected 1 entry, got %d", len(ff))
	}
}

func TestFreezeKeys_AddsEntries(t *testing.T) {
	ff := make(FreezeFile)
	result := FreezeKeys(ff, []string{"TOKEN", "SECRET"}, "deploy-bot")
	if len(result) != 2 {
		t.Fatalf("expected 2 frozen keys, got %d", len(result))
	}
	if result["TOKEN"].FrozenBy != "deploy-bot" {
		t.Errorf("expected FrozenBy=deploy-bot")
	}
	if result["TOKEN"].FrozenAt.IsZero() {
		t.Errorf("expected non-zero FrozenAt")
	}
}

func TestFilterFrozen_RemovesFrozenKeys(t *testing.T) {
	ff := FreezeFile{
		"DB_PASS": {Key: "DB_PASS", FrozenAt: time.Now(), FrozenBy: "alice"},
	}
	secrets := map[string]string{
		"DB_PASS": "secret",
		"APP_ENV": "production",
	}
	out, dropped := FilterFrozen(secrets, ff)
	if len(out) != 1 {
		t.Errorf("expected 1 key after filter, got %d", len(out))
	}
	if _, ok := out["APP_ENV"]; !ok {
		t.Errorf("expected APP_ENV to remain")
	}
	if len(dropped) != 1 || dropped[0] != "DB_PASS" {
		t.Errorf("expected dropped=[DB_PASS], got %v", dropped)
	}
}

func TestFilterFrozen_EmptyFreezeFile(t *testing.T) {
	ff := make(FreezeFile)
	secrets := map[string]string{"A": "1", "B": "2"}
	out, dropped := FilterFrozen(secrets, ff)
	if len(out) != 2 {
		t.Errorf("expected all keys to pass through")
	}
	if len(dropped) != 0 {
		t.Errorf("expected no dropped keys")
	}
}
