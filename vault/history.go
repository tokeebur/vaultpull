package vault

import (
	"encoding/json"
	"os"
	"time"
)

// HistoryEntry records a snapshot of secrets at a point in time.
type HistoryEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Secrets   map[string]string `json:"secrets"`
}

// AppendHistory appends a new history entry to the given file.
func AppendHistory(path, source string, secrets map[string]string) error {
	entries, err := LoadHistory(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	entries = append(entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Secrets:   copyMap(secrets),
	})
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadHistory loads all history entries from the given file.
func LoadHistory(path string) ([]HistoryEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// DiffHistory returns keys that changed between two history entries.
func DiffHistory(a, b HistoryEntry) map[string][2]string {
	changes := make(map[string][2]string)
	keys := make(map[string]struct{})
	for k := range a.Secrets {
		keys[k] = struct{}{}
	}
	for k := range b.Secrets {
		keys[k] = struct{}{}
	}
	for k := range keys {
		av, bv := a.Secrets[k], b.Secrets[k]
		if av != bv {
			changes[k] = [2]string{av, bv}
		}
	}
	return changes
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
