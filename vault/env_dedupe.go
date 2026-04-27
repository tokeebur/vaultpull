package vault

import (
	"fmt"
	"sort"
	"strings"
)

// DedupeResult holds information about a single deduplication action.
type DedupeResult struct {
	Key      string
	Kept     string
	Dropped  []string
	Strategy string
}

// DedupeSecrets removes duplicate values from a secrets map, keeping one
// canonical key per value according to the given strategy.
// Strategies: "first" (keep shortest key), "last" (keep longest key), "alpha" (keep alphabetically first key).
func DedupeSecrets(secrets map[string]string, strategy string) (map[string]string, []DedupeResult, error) {
	if strategy == "" {
		strategy = "alpha"
	}
	switch strategy {
	case "first", "last", "alpha":
	default:
		return nil, nil, fmt.Errorf("unknown deduplication strategy %q: must be first, last, or alpha", strategy)
	}

	// Group keys by their value.
	valueToKeys := make(map[string][]string)
	for k, v := range secrets {
		valueToKeys[v] = append(valueToKeys[v], k)
	}

	result := make(map[string]string, len(secrets))
	var report []DedupeResult

	for v, keys := range valueToKeys {
		if len(keys) == 1 {
			result[keys[0]] = v
			continue
		}
		sort.Strings(keys)
		var kept string
		switch strategy {
		case "alpha":
			kept = keys[0]
		case "first":
			kept = keys[0]
			for _, k := range keys[1:] {
				if len(k) < len(kept) {
					kept = k
				}
			}
		case "last":
			kept = keys[0]
			for _, k := range keys[1:] {
				if len(k) > len(kept) {
					kept = k
				}
			}
		}
		dropped := make([]string, 0, len(keys)-1)
		for _, k := range keys {
			if k != kept {
				dropped = append(dropped, k)
			}
		}
		result[kept] = v
		report = append(report, DedupeResult{Key: kept, Kept: v, Dropped: dropped, Strategy: strategy})
	}

	sort.Slice(report, func(i, j int) bool { return report[i].Key < report[j].Key })
	return result, report, nil
}

// FormatDedupeReport formats deduplication results as a human-readable string.
func FormatDedupeReport(results []DedupeResult) string {
	if len(results) == 0 {
		return "no duplicate values found"
	}
	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "kept %s (dropped: %s)\n", r.Key, strings.Join(r.Dropped, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}
