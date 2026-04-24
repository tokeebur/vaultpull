package vault

import (
	"fmt"
	"sort"
	"strings"
)

// Scope represents a named set of key filters applied to secrets.
type Scope struct {
	Name    string
	Include []string // glob patterns
	Exclude []string // glob patterns
}

// ApplyScope filters secrets to only those matching the scope's include/exclude rules.
// Include patterns are evaluated first; exclude patterns remove keys from the result.
func ApplyScope(secrets map[string]string, scope Scope) (map[string]string, error) {
	result := make(map[string]string)
	for k, v := range secrets {
		included := len(scope.Include) == 0
		for _, pat := range scope.Include {
			if matchesPattern(pat, k) {
				included = true
				break
			}
		}
		if !included {
			continue
		}
		excluded := false
		for _, pat := range scope.Exclude {
			if matchesPattern(pat, k) {
				excluded = true
				break
			}
		}
		if !excluded {
			result[k] = v
		}
	}
	return result, nil
}

// FormatScopeReport returns a human-readable summary of keys included by a scope.
func FormatScopeReport(name string, keys []string) string {
	if len(keys) == 0 {
		return fmt.Sprintf("scope %q: no keys matched", name)
	}
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	var sb strings.Builder
	fmt.Fprintf(&sb, "scope %q matched %d key(s):\n", name, len(sorted))
	for _, k := range sorted {
		fmt.Fprintf(&sb, "  %s\n", k)
	}
	return strings.TrimRight(sb.String(), "\n")
}
