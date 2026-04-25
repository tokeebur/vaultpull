package vault

import (
	"fmt"
	"sort"
	"strings"
)

// PrefixAction describes what operation to perform on secret keys.
type PrefixAction string

const (
	PrefixActionAdd    PrefixAction = "add"
	PrefixActionRemove PrefixAction = "remove"
	PrefixActionRename PrefixAction = "rename"
)

// PrefixResult holds the outcome of a single key transformation.
type PrefixResult struct {
	OldKey string
	NewKey string
	Action PrefixAction
}

// ApplyPrefix adds or removes a prefix from all matching secret keys.
// action: "add" prepends prefix, "remove" strips prefix (skips non-matching keys).
// Returns a new map and a report of changes.
func ApplyPrefix(secrets map[string]string, prefix string, action PrefixAction) (map[string]string, []PrefixResult, error) {
	if prefix == "" {
		return nil, nil, fmt.Errorf("prefix must not be empty")
	}
	if action != PrefixActionAdd && action != PrefixActionRemove {
		return nil, nil, fmt.Errorf("unknown prefix action: %s", action)
	}

	out := make(map[string]string, len(secrets))
	var results []PrefixResult

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := secrets[k]
		switch action {
		case PrefixActionAdd:
			newKey := prefix + k
			out[newKey] = v
			results = append(results, PrefixResult{OldKey: k, NewKey: newKey, Action: PrefixActionAdd})
		case PrefixActionRemove:
			if strings.HasPrefix(k, prefix) {
				newKey := strings.TrimPrefix(k, prefix)
				out[newKey] = v
				results = append(results, PrefixResult{OldKey: k, NewKey: newKey, Action: PrefixActionRemove})
			} else {
				out[k] = v
			}
		}
	}
	return out, results, nil
}

// FormatPrefixReport returns a human-readable summary of prefix changes.
func FormatPrefixReport(results []PrefixResult) string {
	if len(results) == 0 {
		return "no keys changed"
	}
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("[%s] %s -> %s\n", r.Action, r.OldKey, r.NewKey))
	}
	return strings.TrimRight(sb.String(), "\n")
}
