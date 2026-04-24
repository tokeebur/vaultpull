package vault

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// PlaceholderRule defines a key pattern and the placeholder text to use.
type PlaceholderRule struct {
	Pattern     string
	Placeholder string
}

// DefaultPlaceholderRules returns a set of sensible default rules.
func DefaultPlaceholderRules() []PlaceholderRule {
	return []PlaceholderRule{
		{Pattern: "*_URL", Placeholder: "<url>"},
		{Pattern: "*_HOST", Placeholder: "<host>"},
		{Pattern: "*_PORT", Placeholder: "<port>"},
		{Pattern: "*_KEY", Placeholder: "<secret>"},
		{Pattern: "*_SECRET", Placeholder: "<secret>"},
		{Pattern: "*_PASSWORD", Placeholder: "<secret>"},
		{Pattern: "*_TOKEN", Placeholder: "<secret>"},
	}
}

// ApplyPlaceholders replaces values in secrets with placeholder text for keys
// matching the given rules. Keys not matching any rule retain their original
// value. Rules are evaluated in order; the last match wins.
func ApplyPlaceholders(secrets map[string]string, rules []PlaceholderRule) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, rule := range rules {
		pattern, err := globToRegexp(rule.Pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", rule.Pattern, err)
		}
		for k := range out {
			if pattern.MatchString(k) {
				out[k] = rule.Placeholder
			}
		}
	}
	return out, nil
}

// ListPlaceholderKeys returns sorted keys whose values were replaced.
func ListPlaceholderKeys(original, placeholdered map[string]string) []string {
	var keys []string
	for k, v := range placeholdered {
		if orig, ok := original[k]; ok && orig != v {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

func globToRegexp(pattern string) (*regexp.Regexp, error) {
	escaped := regexp.QuoteMeta(pattern)
	regexStr := "^" + strings.ReplaceAll(escaped, `\*`, `.*`) + "$"
	return regexp.Compile(regexStr)
}
