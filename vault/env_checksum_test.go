package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComputeChecksum_Deterministic(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	a := ComputeChecksum(secrets)
	b := ComputeChecksum(secrets)
	if a.Digest != b.Digest {
		t.Fatalf("expected deterministic digest, got %q vs %q", a.Digest, b.Digest)
	}
	if len(a.Keys) != 2 {
		t.Fatalf("expected 2 per-key digests, got %d", len(a.Keys))
	}
}

func TestComputeChecksum_DifferentValues(t *testing.T) {
	a := ComputeChecksum(map[string]string{"KEY": "value1"})
	b := ComputeChecksum(map[string]string{"KEY": "value2"})
	if a.Digest == b.Digest {
		t.Fatal("expected different digests for different values")
	}
}

func TestComputeChecksum_EmptyMap(t *testing.T) {
	c := ComputeChecksum(map[string]string{})
	if c.Digest == "" {
		t.Fatal("expected non-empty digest for empty map")
	}
	if len(c.Keys) != 0 {
		t.Fatalf("expected no per-key digests, got %d", len(c.Keys))
	}
}

func TestSaveAndLoadChecksumFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "checksum.json")
	secrets := map[string]string{"A": "1", "B": "2"}
	orig := ComputeChecksum(secrets)

	if err := SaveChecksumFile(path, orig); err != nil {
		t.Fatalf("SaveChecksumFile: %v", err)
	}
	loaded, err := LoadChecksumFile(path)
	if err != nil {
		t.Fatalf("LoadChecksumFile: %v", err)
	}
	if loaded.Digest != orig.Digest {
		t.Fatalf("digest mismatch: %q vs %q", loaded.Digest, orig.Digest)
	}
}

func TestLoadChecksumFile_NotExist(t *testing.T) {
	c, err := LoadChecksumFile("/nonexistent/path/checksum.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if c.Digest != "" {
		t.Fatalf("expected empty digest, got %q", c.Digest)
	}
}

func TestVerifyChecksum_Match(t *testing.T) {
	secrets := map[string]string{"X": "hello"}
	stored := ComputeChecksum(secrets)
	if !VerifyChecksum(secrets, stored) {
		t.Fatal("expected checksum to match")
	}
}

func TestVerifyChecksum_Mismatch(t *testing.T) {
	stored := ComputeChecksum(map[string]string{"X": "hello"})
	if VerifyChecksum(map[string]string{"X": "world"}, stored) {
		t.Fatal("expected checksum mismatch")
	}
}

func TestChangedKeys_DetectsChanges(t *testing.T) {
	stored := ComputeChecksum(map[string]string{"A": "1", "B": "2", "C": "3"})
	updated := map[string]string{"A": "1", "B": "changed", "C": "3"}
	changed := ChangedKeys(updated, stored)
	if len(changed) != 1 || changed[0] != "B" {
		t.Fatalf("expected [B], got %v", changed)
	}
}

func TestSaveChecksumFile_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "checksum.json")
	if err := SaveChecksumFile(path, ComputeChecksum(map[string]string{})); err != nil {
		t.Fatalf("SaveChecksumFile: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
