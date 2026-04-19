package vault

import (
	"fmt"
	"strings"
)

// MaskRule defines how a key should be masked.
type MaskRule struct {
	Key      string
	ShowLast int // number of trailing chars to reveal
}

// DefaultMaskRules returns common sensitive key patterns.
func DefaultMaskRules() []MaskRule {
	return []MaskRule{
		{Key: "PASSWORD", ShowLast: 0},
		{Key: "SECRET", ShowLast: 0},
		{Key: "TOKEN", ShowLast: 4},
		{Key: "API_KEY", ShowLast: 4},
		{Key: "PRIVATE_KEY", ShowLast: 0},
	}
}

// MaskSecrets returns a copy of secrets with sensitive values masked.
func MaskSecrets(secrets map[string]string, rules []MaskRule) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = maskValue(k, v, rules)
	}
	return result
}

func maskValue(key, value string, rules []MaskRule) string {
	upperKey := strings.ToUpper(key)
	for _, rule := range rules {
		if strings.Contains(upperKey, rule.Key) {
			if rule.ShowLast <= 0 || len(value) <= rule.ShowLast {
				return "****"
			}
			return fmt.Sprintf("****%s", value[len(value)-rule.ShowLast:])
		}
	}
	return value
}
