package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadRefFile reads a ref definition file.
// Each non-comment line has the form:  KEY  SOURCE  [note...]
func LoadRefFile(path string) ([]RefEntry, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []RefEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid ref line: %q", line)
		}
		e := RefEntry{Key: parts[0], Source: parts[1]}
		if len(parts) > 2 {
			e.Note = strings.Join(parts[2:], " ")
		}
		entries = append(entries, e)
	}
	return entries, scanner.Err()
}

// SaveRefFile writes ref entries to path.
func SaveRefFile(path string, refs []RefEntry) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "# vaultpull ref file — KEY SOURCE [note]")
	for _, r := range refs {
		line := fmt.Sprintf("%s %s", r.Key, r.Source)
		if r.Note != "" {
			line += " " + r.Note
		}
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
