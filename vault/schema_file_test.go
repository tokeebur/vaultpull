package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveSchemaFile_Content(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.txt")
	schema := SchemaFile{Fields: []FieldSchema{
		{Key: "TOKEN", Type: "string", Required: true, Pattern: "-"},
	}}
	if err := SaveSchemaFile(path, schema); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "TOKEN") {
		t.Error("expected TOKEN in schema file")
	}
	if !strings.Contains(string(data), "string") {
		t.Error("expected type string in schema file")
	}
}

func TestLoadSchemaFile_RequiredParsed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.txt")
	os.WriteFile(path, []byte("SECRET string true -\n"), 0644)
	loaded, err := LoadSchemaFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.Fields[0].Required {
		t.Error("expected Required=true")
	}
}

func TestLoadSchemaFile_PatternParsed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.txt")
	os.WriteFile(path, []byte(`PORT int true ^\d+$`+"\n"), 0644)
	loaded, err := LoadSchemaFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Fields[0].Pattern != `^\d+$` {
		t.Errorf("unexpected pattern: %s", loaded.Fields[0].Pattern)
	}
}

func TestSaveSchemaFile_Permission(t *testing.T) {
	err := SaveSchemaFile("/nonexistent/dir/schema.txt", SchemaFile{})
	if err == nil {
		t.Fatal("expected error writing to invalid path")
	}
}
