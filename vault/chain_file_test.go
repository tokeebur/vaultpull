package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadChainFile_NotExist(t *testing.T) {
	_, err := LoadChainFile("/nonexistent/chain.conf")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAndLoadChainFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chain.conf")

	cc := &ChainConfig{
		Sources: []ChainSource{
			{Name: "base", Path: "secret/base"},
			{Name: "prod", Path: "secret/prod"},
		},
	}
	if err := SaveChainFile(path, cc); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadChainFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(loaded.Sources))
	}
	if loaded.Sources[0].Name != "base" || loaded.Sources[1].Path != "secret/prod" {
		t.Error("unexpected source content")
	}
}

func TestLoadChainFile_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chain.conf")
	content := "# comment\nbase=secret/base\n"
	os.WriteFile(path, []byte(content), 0644)
	cc, err := LoadChainFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cc.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(cc.Sources))
	}
}

func TestLoadChainFile_InvalidLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chain.conf")
	os.WriteFile(path, []byte("badline\n"), 0644)
	_, err := LoadChainFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}
