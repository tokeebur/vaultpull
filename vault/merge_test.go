package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeSecrets_NewFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")

	remote := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := MergeSecrets(out, remote, MergeStrategyKeepRemote)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(res.Added))
	}
	if len(res.Updated) != 0 || len(res.Skipped) != 0 {
		t.Errorf("expected no updates or skips")
	}
	if _, err := os.Stat(out); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestMergeSecrets_KeepLocal(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	_ = os.WriteFile(out, []byte("FOO=local\nEXISTING=yes\n"), 0600)

	remote := map[string]string{"FOO": "remote", "NEW": "value"}
	res, err := MergeSecrets(out, remote, MergeStrategyKeepLocal)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected FOO to be skipped, got %v", res.Skipped)
	}
	if len(res.Added) != 1 || res.Added[0] != "NEW" {
		t.Errorf("expected NEW to be added, got %v", res.Added)
	}
	if res.Final["FOO"] != "local" {
		t.Errorf("expected FOO=local, got %s", res.Final["FOO"])
	}
}

func TestMergeSecrets_KeepRemote(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	_ = os.WriteFile(out, []byte("FOO=local\n"), 0600)

	remote := map[string]string{"FOO": "remote"}
	res, err := MergeSecrets(out, remote, MergeStrategyKeepRemote)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Updated) != 1 || res.Updated[0] != "FOO" {
		t.Errorf("expected FOO to be updated, got %v", res.Updated)
	}
	if res.Final["FOO"] != "remote" {
		t.Errorf("expected FOO=remote, got %s", res.Final["FOO"])
	}
}

func TestMergeSecrets_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env")
	remote := map[string]string{"SECRET": "value"}
	_, err := MergeSecrets(out, remote, MergeStrategyKeepRemote)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
