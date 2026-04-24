package vault

import (
	"fmt"
	"sort"
	"strings"
)

// RequiredRule defines a key that must be present and optionally non-empty.
type RequiredRule struct {
	Key      string
	NonEmpty bool
}

// RequiredViolation describes a single requirement that was not satisfied.
type RequiredViolation struct {
	Key    string
	Reason string
}

// CheckRequired verifies that all required keys are present in secrets.
// If NonEmpty is true for a rule, the value must also be non-empty.
func CheckRequired(secrets map[string]string, rules []RequiredRule) []RequiredViolation {
	var violations []RequiredViolation
	for _, rule := range rules {
		val, ok := secrets[rule.Key]
		if !ok {
			violations = append(violations, RequiredViolation{
				Key:    rule.Key,
				Reason: "missing",
			})
			continue
		}
		if rule.NonEmpty && strings.TrimSpace(val) == "" {
			violations = append(violations, RequiredViolation{
				Key:    rule.Key,
				Reason: "empty value",
			})
		}
	}
	return violations
}

// FormatRequiredReport returns a human-readable summary of violations.
func FormatRequiredReport(violations []RequiredViolation) string {
	if len(violations) == 0 {
		return "all required keys satisfied"
	}
	sorted := make([]RequiredViolation, len(violations))
	copy(sorted, violations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d required key violation(s):\n", len(sorted)))
	for _, v := range sorted {
		sb.WriteString(fmt.Sprintf("  %-30s %s\n", v.Key, v.Reason))
	}
	return strings.TrimRight(sb.String(), "\n")
}
