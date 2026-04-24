package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadDeprecationFile reads a deprecation rules file.
// Each non-comment line has the format:
//   OLD_KEY=NEW_KEY:reason
// NEW_KEY and reason are optional.
func LoadDeprecationFile(path string) (map[string]DeprecationEntry, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return map[string]DeprecationEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rules := map[string]DeprecationEntry{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		entry := DeprecationEntry{Key: key}
		if len(parts) == 2 {
			rhs := strings.TrimSpace(parts[1])
			if idx := strings.Index(rhs, ":"); idx >= 0 {
				entry.Replacement = strings.TrimSpace(rhs[:idx])
				entry.Reason = strings.TrimSpace(rhs[idx+1:])
			} else {
				entry.Replacement = rhs
			}
		}
		rules[key] = entry
	}
	return rules, scanner.Err()
}

// SaveDeprecationFile writes deprecation rules to a file.
func SaveDeprecationFile(path string, rules map[string]DeprecationEntry) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "# vaultpull deprecation rules: OLD_KEY=NEW_KEY:reason")
	for key, e := range rules {
		line := key
		if e.Replacement != "" || e.Reason != "" {
			line += "=" + e.Replacement + ":" + e.Reason
		}
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
