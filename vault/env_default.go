package vault

import (
	"fmt"
	"strings"
)

// DefaultRule defines a fallback value for a secret key if it is missing or empty.
type DefaultRule struct {
	Key      string
	Value    string
	OnEmpty  bool // apply default even when key exists but value is empty
}

// ApplyDefaults fills in missing or empty secret values based on the provided rules.
// Returns a new map with defaults applied and a report of which keys were defaulted.
func ApplyDefaults(secrets map[string]string, rules []DefaultRule) (map[string]string, []string) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var applied []string
	for _, rule := range rules {
		if rule.Key == "" {
			continue
		}
		existing, exists := out[rule.Key]
		if !exists || (rule.OnEmpty && strings.TrimSpace(existing) == "") {
			out[rule.Key] = rule.Value
			applied = append(applied, rule.Key)
		}
	}
	return out, applied
}

// FormatDefaultReport returns a human-readable summary of applied defaults.
func FormatDefaultReport(applied []string) string {
	if len(applied) == 0 {
		return "No defaults applied."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d default(s) applied:\n", len(applied)))
	for _, key := range applied {
		sb.WriteString(fmt.Sprintf("  - %s\n", key))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// ParseDefaultRules parses lines of the form "KEY=value" or "KEY?=value" into DefaultRules.
// The "?=" form sets OnEmpty=true.
func ParseDefaultRules(lines []string) ([]DefaultRule, error) {
	var rules []DefaultRule
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		onEmpty := false
		sep := "="
		if strings.Contains(line, "?=") {
			sep = "?="
			onEmpty = true
		}
		parts := strings.SplitN(line, sep, 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid default rule %q", i+1, line)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("line %d: empty key in default rule", i+1)
		}
		rules = append(rules, DefaultRule{Key: key, Value: val, OnEmpty: onEmpty})
	}
	return rules, nil
}
