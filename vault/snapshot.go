package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of secrets.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
}

// SaveSnapshot writes a snapshot of secrets to a JSON file in dir.
func SaveSnapshot(dir, vaultPath string, secrets map[string]string) (string, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create snapshot dir: %w", err)
	}
	s := Snapshot{
		Timestamp: time.Now().UTC(),
		Path:      vaultPath,
		Secrets:   secrets,
	}
	filename := fmt.Sprintf("snapshot_%d.json", s.Timestamp.UnixNano())
	out := filepath.Join(dir, filename)
	f, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return "", fmt.Errorf("open snapshot file: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(s); err != nil {
		return "", fmt.Errorf("encode snapshot: %w", err)
	}
	return out, nil
}

// LoadSnapshot reads a snapshot from a JSON file.
func LoadSnapshot(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("decode snapshot: %w", err)
	}
	return &s, nil
}
