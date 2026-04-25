package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAnnotationFile_NotExist(t *testing.T) {
	am, err := LoadAnnotationFile("/tmp/vaultpull_nonexistent_annotations.txt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(am) != 0 {
		t.Errorf("expected empty map, got %v", am)
	}
}

func TestSaveAndLoadAnnotationFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "annotations.txt")

	am := AnnotationMap{
		"DB_PASSWORD": {"sensitive credential", "rotated quarterly"},
		"API_KEY":     {"third-party service key"},
	}

	if err := SaveAnnotationFile(path, am); err != nil {
		t.Fatalf("SaveAnnotationFile: %v", err)
	}

	loaded, err := LoadAnnotationFile(path)
	if err != nil {
		t.Fatalf("LoadAnnotationFile: %v", err)
	}

	for key, notes := range am {
		got, ok := loaded[key]
		if !ok {
			t.Errorf("missing key %s", key)
			continue
		}
		if len(got) != len(notes) {
			t.Errorf("key %s: expected %d notes, got %d", key, len(notes), len(got))
		}
	}
}

func TestLoadAnnotationFile_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "annotations.txt")

	content := "# this is a comment\nDB_HOST internal database host\n\nAPI_KEY service key\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	am, err := LoadAnnotationFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(am) != 2 {
		t.Errorf("expected 2 entries, got %d", len(am))
	}
}

func TestAddAnnotation_AppendsNote(t *testing.T) {
	am := make(AnnotationMap)
	am = AddAnnotation(am, "SECRET_KEY", "first note")
	am = AddAnnotation(am, "SECRET_KEY", "second note")

	if len(am["SECRET_KEY"]) != 2 {
		t.Errorf("expected 2 notes, got %d", len(am["SECRET_KEY"]))
	}
}

func TestFormatAnnotations_Output(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS": "secret",
		"API_KEY": "key123",
	}
	am := AnnotationMap{
		"DB_PASS": {"sensitive"},
	}

	out := FormatAnnotations(secrets, am)
	if out == "" {
		t.Error("expected non-empty output")
	}
	if !containsStr(out, "DB_PASS") {
		t.Error("expected DB_PASS in output")
	}
	if containsStr(out, "API_KEY") {
		t.Error("API_KEY has no annotation, should not appear")
	}
}
