package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Checksum represents a SHA-256 digest of a secrets map.
type Checksum struct {
	Digest string            `json:"digest"`
	Keys   map[string]string `json:"keys"` // per-key digests
}

// ComputeChecksum returns a Checksum for the given secrets map.
// Each key is hashed individually, and a combined digest is produced
// from the sorted key=hash pairs for determinism.
func ComputeChecksum(secrets map[string]string) Checksum {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	perKey := make(map[string]string, len(secrets))
	var combined strings.Builder
	for _, k := range keys {
		h := sha256.Sum256([]byte(secrets[k]))
		d := hex.EncodeToString(h[:])
		perKey[k] = d
		fmt.Fprintf(&combined, "%s=%s\n", k, d)
	}

	total := sha256.Sum256([]byte(combined.String()))
	return Checksum{
		Digest: hex.EncodeToString(total[:]),
		Keys:   perKey,
	}
}

// SaveChecksumFile writes a Checksum to path as JSON.
func SaveChecksumFile(path string, c Checksum) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("checksum marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// LoadChecksumFile reads a Checksum from a JSON file.
// Returns an empty Checksum (zero Digest) if the file does not exist.
func LoadChecksumFile(path string) (Checksum, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Checksum{}, nil
	}
	if err != nil {
		return Checksum{}, fmt.Errorf("checksum read: %w", err)
	}
	var c Checksum
	if err := json.Unmarshal(data, &c); err != nil {
		return Checksum{}, fmt.Errorf("checksum parse: %w", err)
	}
	return c, nil
}

// VerifyChecksum compares a freshly computed checksum against a stored one.
// Returns true when both digests match.
func VerifyChecksum(secrets map[string]string, stored Checksum) bool {
	fresh := ComputeChecksum(secrets)
	return fresh.Digest == stored.Digest
}

// ChangedKeys returns keys whose per-key digest differs between stored and fresh.
func ChangedKeys(secrets map[string]string, stored Checksum) []string {
	fresh := ComputeChecksum(secrets)
	var changed []string
	for k, d := range fresh.Keys {
		if stored.Keys[k] != d {
			changed = append(changed, k)
		}
	}
	sort.Strings(changed)
	return changed
}
