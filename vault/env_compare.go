package vault

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two secret maps.
type CompareResult struct {
	OnlyInA  []string
	OnlyInB  []string
	Differ   []string
	Identical []string
}

// CompareSecrets compares two secret maps and returns a CompareResult.
func CompareSecrets(a, b map[string]string) CompareResult {
	result := CompareResult{}
	keys := map[string]struct{}{}
	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}
	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	for _, k := range sorted {
		av, aok := a[k]
		bv, bok := b[k]
		switch {
		case aok && !bok:
			result.OnlyInA = append(result.OnlyInA, k)
		case !aok && bok:
			result.OnlyInB = append(result.OnlyInB, k)
		case av == bv:
			result.Identical = append(result.Identical, k)
		default:
			result.Differ = append(result.Differ, k)
		}
	}
	return result
}

// FormatCompareResult returns a human-readable summary of a CompareResult.
func FormatCompareResult(r CompareResult, labelA, labelB string) string {
	var sb strings.Builder
	for _, k := range r.OnlyInA {
		fmt.Fprintf(&sb, "< %s (only in %s)\n", k, labelA)
	}
	for _, k := range r.OnlyInB {
		fmt.Fprintf(&sb, "> %s (only in %s)\n", k, labelB)
	}
	for _, k := range r.Differ {
		fmt.Fprintf(&sb, "~ %s (differs)\n", k)
	}
	for _, k := range r.Identical {
		fmt.Fprintf(&sb, "= %s (identical)\n", k)
	}
	return sb.String()
}
