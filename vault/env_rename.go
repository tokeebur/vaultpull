package vault

import (
	"fmt"
	"strings"
)

// RenameRule defines a single key rename operation.
type RenameRule struct {
	From string
	To   string
}

// RenameSecrets applies rename rules to a secrets map, returning a new map.
// If a source key does not exist, it is skipped unless strict is true.
func RenameSecrets(secrets map[string]string, rules []RenameRule, strict bool) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, r := range rules {
		val, ok := out[r.From]
		if !ok {
			if strict {
				return nil, fmt.Errorf("rename: source key %q not found", r.From)
			}
			continue
		}
		delete(out, r.From)
		out[r.To] = val
	}
	return out, nil
}

// ParseRenameRules parses lines of the form "FROM=TO" into RenameRule slice.
// Lines starting with '#' or empty lines are skipped.
func ParseRenameRules(lines []string) ([]RenameRule, error) {
	var rules []RenameRule
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("rename: invalid rule on line %d: %q", i+1, line)
		}
		rules = append(rules, RenameRule{From: strings.TrimSpace(parts[0]), To: strings.TrimSpace(parts[1])})
	}
	return rules, nil
}
