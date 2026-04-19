package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProfileFile_NotExist(t *testing.T) {
	profiles, err := LoadProfileFile("/tmp/nonexistent_profile_xyz.ini")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(profiles) != 0 {
		t.Fatalf("expected empty profiles")
	}
}

func TestSaveAndLoadProfileFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.ini")

	profiles := map[string]Profile{
		"dev": {Name: "dev", Values: map[string]string{"DB_HOST": "localhost", "DEBUG": "true"}},
		"prod": {Name: "prod", Values: map[string]string{"DB_HOST": "prod.db", "DEBUG": "false"}},
	}
	if err := SaveProfileFile(path, profiles); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadProfileFile(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded["dev"].Values["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", loaded["dev"].Values["DB_HOST"])
	}
	if loaded["prod"].Values["DEBUG"] != "false" {
		t.Errorf("expected false, got %s", loaded["prod"].Values["DEBUG"])
	}
}

func TestLoadProfileFile_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.ini")
	content := "# comment\n[staging]\n# another comment\nKEY=val\n"
	os.WriteFile(path, []byte(content), 0600)
	profiles, err := LoadProfileFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profiles["staging"].Values["KEY"] != "val" {
		t.Errorf("expected val")
	}
}

func TestApplyProfile_OverridesSecrets(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	p := Profile{Name: "test", Values: map[string]string{"B": "99", "C": "3"}}
	result := ApplyProfile(secrets, p)
	if result["A"] != "1" {
		t.Errorf("expected A=1")
	}
	if result["B"] != "99" {
		t.Errorf("expected B=99")
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3")
	}
}

func TestApplyProfile_DoesNotMutateOriginal(t *testing.T) {
	secrets := map[string]string{"X": "orig"}
	p := Profile{Name: "p", Values: map[string]string{"X": "override"}}
	ApplyProfile(secrets, p)
	if secrets["X"] != "orig" {
		t.Errorf("original map was mutated")
	}
}
