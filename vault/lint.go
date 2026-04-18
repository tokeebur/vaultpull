package vault

import (
	"fmt"
	"strings"
)

// LintRule represents a single linting rule.
type LintRule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// LintResult holds the result of a lint check.
type LintResult struct {
	Key     string
	Rule    string
	Message string
}

// DefaultRules returns the built-in lint rules.
func DefaultRules() []LintRule {
	return []LintRule{
		{
			Name:    "no-empty-value",
			Message: "value must not be empty",
			Check:   func(_, v string) bool { return strings.TrimSpace(v) == "" },
		},
		{
			Name:    "no-lowercase-key",
			Message: "key should be uppercase",
			Check:   func(k, _ string) bool { return k != strings.ToUpper(k) },
		},
		{
			Name:    "no-space-in-key",
			Message: "key must not contain spaces",
			Check:   func(k, _ string) bool { return strings.Contains(k, " ") },
		},
		{
			Name:    "no-secret-in-key-name",
			Message: "key name should not contain 'SECRET' literally (use a descriptive name)",
			Check:   func(k, _ string) bool { return strings.EqualFold(k, "SECRET") },
		},
	}
}

// LintSecrets runs lint rules against a map of secrets and returns any violations.
func LintSecrets(secrets map[string]string, rules []LintRule) []LintResult {
	var results []LintResult
	for k, v := range secrets {
		for _, rule := range rules {
			if rule.Check(k, v) {
				results = append(results, LintResult{
					Key:     k,
					Rule:    rule.Name,
					Message: fmt.Sprintf("%s: %s", k, rule.Message),
				})
			}
		}
	}
	return results
}
