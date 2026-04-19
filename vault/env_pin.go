package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PinEntry holds a pinned key-value pair with an optional comment.
type PinEntry struct {
	Key     string
	Value   string
	Comment string
}

// PinFile maps keys to PinEntry.
type PinFile map[string]PinEntry

// LoadPinFile reads a pin file from disk.
func LoadPinFile(path string) (PinFile, error) {
	pf := make(PinFile)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return pf, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		rest := parts[1]
		comment := ""
		if idx := strings.Index(rest, " #"); idx != -1 {
			comment = strings.TrimSpace(rest[idx+2:])
			rest = rest[:idx]
		}
		pf[key] = PinEntry{Key: key, Value: strings.TrimSpace(rest), Comment: comment}
	}
	return pf, scanner.Err()
}

// SavePinFile writes pin entries to disk.
func SavePinFile(path string, pf PinFile) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, e := range pf {
		line := fmt.Sprintf("%s=%s", e.Key, e.Value)
		if e.Comment != "" {
			line += " # " + e.Comment
		}
		fmt.Fprintln(f, line)
	}
	return nil
}

// ApplyPins overrides secrets with pinned values.
func ApplyPins(secrets map[string]string, pf PinFile) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for key, entry := range pf {
		out[key] = entry.Value
	}
	return out
}

// PinnedKeys returns the list of keys currently pinned.
func PinnedKeys(pf PinFile) []string {
	keys := make([]string, 0, len(pf))
	for k := range pf {
		keys = append(keys, k)
	}
	return keys
}
