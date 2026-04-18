package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenderTemplate_Success(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "app.env.tmpl")
	destPath := filepath.Join(dir, "app.env")

	tmplContent := "DB_HOST={{.db_host}}\nDB_PASS={{.db_pass}}\n"
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatal(err)
	}

	secrets := map[string]interface{}{
		"db_host": "localhost",
		"db_pass": "s3cr3t",
	}

	if err := RenderTemplate(tmplPath, destPath, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatal(err)
	}

	expected := "DB_HOST=localhost\nDB_PASS=s3cr3t\n"
	if string(data) != expected {
		t.Errorf("got %q, want %q", string(data), expected)
	}
}

func TestRenderTemplate_MissingKey(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "app.env.tmpl")
	destPath := filepath.Join(dir, "app.env")

	if err := os.WriteFile(tmplPath, []byte("VAL={{.missing_key}}"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := RenderTemplate(tmplPath, destPath, map[string]interface{}{}); err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestRenderTemplate_TemplateNotExist(t *testing.T) {
	err := RenderTemplate("/nonexistent/path.tmpl", "/tmp/out", nil)
	if err == nil {
		t.Error("expected error for missing template file")
	}
}

func TestRenderTemplate_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "app.env.tmpl")
	destPath := filepath.Join(dir, "app.env")

	if err := os.WriteFile(tmplPath, []byte("KEY={{.key}}"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := RenderTemplate(tmplPath, destPath, map[string]interface{}{"key": "val"}); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(destPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}

func TestRenderTemplate_DestNotCreatedOnError(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "app.env.tmpl")
	destPath := filepath.Join(dir, "app.env")

	if err := os.WriteFile(tmplPath, []byte("VAL={{.missing}}"), 0644); err != nil {
		t.Fatal(err)
	}

	// Rendering should fail due to missing key; dest file should not be created.
	_ = RenderTemplate(tmplPath, destPath, map[string]interface{}{})

	if _, err := os.Stat(destPath); err == nil {
		t.Error("expected dest file to not exist after failed render")
	}
}
