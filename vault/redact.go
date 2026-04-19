package vault

import (
	"regexp"
	"strings"
)

// RedactRule defines a pattern and replacement for redacting secret values.
type RedactRule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

var defaultRedactRules = []RedactRule{
	{Pattern: regexp.MustCompile(`(?i)password`), Replacement: "***"},
	{Pattern: regexp.MustCompile(`(?i)secret`), Replacement: "***"},
	{Pattern: regexp.MustCompile(`(?i)token`), Replacement: "***"},
	{Pattern: regexp.MustCompile(`(?i)api_key`), Replacement: "***"},
	{Pattern: regexp.MustCompile(`(?i)private_key`), Replacement: "***"},
}

// RedactSecrets returns a copy of secrets with sensitive values replaced.
// Keys matching any rule pattern have their values replaced with the rule's replacement.
func RedactSecrets(secrets map[string]string, rules []RedactRule) map[string]string {
	if rules == nil {
		rules = defaultRedactRules
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = redactValue(k, v, rules)
	}
	return out
}

func redactValue(key, value string, rules []RedactRule) string {
	for _, r := range rules {
		if r.Pattern.MatchString(key) {
			return r.Replacement
		}
	}
	return value
}

// RedactString replaces occurrences of any secret value found in src with "***".
// Useful for sanitising log lines.
func RedactString(src string, secrets map[string]string) string {
	for _, v := range secrets {
		if v == "" {
			continue
		}
		src = strings.ReplaceAll(src, v, "***")
	}
	return src
}
