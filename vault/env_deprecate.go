package vault

import (
	"fmt"
	"sort"
	"strings"
)

// DeprecationEntry holds metadata about a deprecated secret key.
type DeprecationEntry struct {
	Key        string
	Replacement string
	Reason     string
}

// DeprecationReport summarises the result of scanning secrets.
type DeprecationReport struct {
	Deprecated []DeprecationEntry
	Missing    []string // deprecated keys present in secrets
}

// CheckDeprecations scans secrets against a deprecation map and returns a
// report of any deprecated keys that are still in use.
func CheckDeprecations(secrets map[string]string, rules map[string]DeprecationEntry) DeprecationReport {
	report := DeprecationReport{}
	for key := range secrets {
		if entry, ok := rules[key]; ok {
			report.Deprecated = append(report.Deprecated, entry)
		}
	}
	sort.Slice(report.Deprecated, func(i, j int) bool {
		return report.Deprecated[i].Key < report.Deprecated[j].Key
	})
	return report
}

// FormatDeprecationReport returns a human-readable summary of the report.
func FormatDeprecationReport(r DeprecationReport) string {
	if len(r.Deprecated) == 0 {
		return "No deprecated keys found."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d deprecated key(s) found:\n", len(r.Deprecated)))
	for _, e := range r.Deprecated {
		sb.WriteString(fmt.Sprintf("  - %s", e.Key))
		if e.Replacement != "" {
			sb.WriteString(fmt.Sprintf(" → use %s instead", e.Replacement))
		}
		if e.Reason != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", e.Reason))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
