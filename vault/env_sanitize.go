package vault

import (
	"regexp"
	"strings"
)

// SanitizeRule defines a pattern and the action to apply to matching keys.
type SanitizeRule struct {
	Pattern string
	Action  string // "strip_nonprint", "trim", "collapse_space", "remove"
}

// DefaultSanitizeRules returns a baseline set of sanitize rules.
func DefaultSanitizeRules() []SanitizeRule {
	return []SanitizeRule{
		{Pattern: "*", Action: "strip_nonprint"},
		{Pattern: "*", Action: "trim"},
	}
}

// SanitizeSecrets applies sanitize rules to a copy of secrets and returns the
// sanitized map along with a report of keys that were changed.
func SanitizeSecrets(secrets map[string]string, rules []SanitizeRule) (map[string]string, []string) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var changed []string

	for _, rule := range rules {
		re, err := globToRegexp(rule.Pattern)
		if err != nil {
			continue
		}
		pat := regexp.MustCompile(re)
		for k, v := range out {
			if !pat.MatchString(k) {
				continue
			}
			newVal := applySanitizeAction(v, rule.Action)
			if newVal != v {
				out[k] = newVal
				changed = appendUnique(changed, k)
			}
		}
	}

	return out, changed
}

func applySanitizeAction(v, action string) string {
	switch action {
	case "trim":
		return strings.TrimSpace(v)
	case "strip_nonprint":
		return stripNonPrintable(v)
	case "collapse_space":
		space := regexp.MustCompile(`\s+`)
		return space.ReplaceAllString(v, " ")
	case "remove":
		return ""
	}
	return v
}

func stripNonPrintable(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 0x20 || r == '\t' || r == '\n' || r == '\r' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func appendUnique(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}
