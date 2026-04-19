package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateAgainstSchema_AllValid(t *testing.T) {
	schema := SchemaFile{Fields: []FieldSchema{
		{Key: "PORT", Type: "int", Required: true, Pattern: `^\d+$`},
		{Key: "DEBUG", Type: "bool", Required: false},
	}}
	secrets := map[string]string{"PORT": "8080", "DEBUG": "true"}
	violations := ValidateAgainstSchema(secrets, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got: %v", violations)
	}
}

func TestValidateAgainstSchema_MissingRequired(t *testing.T) {
	schema := SchemaFile{Fields: []FieldSchema{
		{Key: "API_KEY", Type: "string", Required: true},
	}}
	violations := ValidateAgainstSchema(map[string]string{}, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got: %v", violations)
	}
}

func TestValidateAgainstSchema_BadInt(t *testing.T) {
	schema := SchemaFile{Fields: []FieldSchema{
		{Key: "PORT", Type: "int", Required: true},
	}}
	violations := ValidateAgainstSchema(map[string]string{"PORT": "abc"}, schema)
	if len(violations) == 0 {
		t.Fatal("expected violation for bad int")
	}
}

func TestValidateAgainstSchema_BadBool(t *testing.T) {
	schema := SchemaFile{Fields: []FieldSchema{
		{Key: "FLAG", Type: "bool"},
	}}
	violations := ValidateAgainstSchema(map[string]string{"FLAG": "yes"}, schema)
	if len(violations) == 0 {
		t.Fatal("expected violation for bad bool")
	}
}

func TestSaveAndLoadSchemaFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.txt")
	schema := SchemaFile{Fields: []FieldSchema{
		{Key: "HOST", Type: "string", Required: true, Pattern: "-"},
		{Key: "PORT", Type: "int", Required: false, Pattern: `^\d+$`},
	}}
	if err := SaveSchemaFile(path, schema); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadSchemaFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(loaded.Fields))
	}
}

func TestLoadSchemaFile_NotExist(t *testing.T) {
	_, err := LoadSchemaFile("/nonexistent/schema.txt")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSchemaFile_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.txt")
	os.WriteFile(path, []byte("# comment\nDB_URL string true -\n"), 0644)
	loaded, err := LoadSchemaFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(loaded.Fields))
	}
}
