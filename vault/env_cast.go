package vault

import (
	"fmt"
	"strconv"
	"strings"
)

// CastType represents a target type for casting a secret value.
type CastType string

const (
	CastString CastType = "string"
	CastInt    CastType = "int"
	CastBool   CastType = "bool"
	CastFloat  CastType = "float"
)

// CastRule defines a key pattern and its target type.
type CastRule struct {
	Pattern string
	Type    CastType
}

// CastResult holds the original and cast value for a key.
type CastResult struct {
	Key      string
	Original string
	Cast     interface{}
	Type     CastType
}

// CastSecrets applies type casting rules to secrets and returns typed results.
func CastSecrets(secrets map[string]string, rules []CastRule) ([]CastResult, error) {
	var results []CastResult
	for _, rule := range rules {
		for k, v := range secrets {
			if !matchesPattern(rule.Pattern, k) {
				continue
			}
			cast, err := castValue(v, rule.Type)
			if err != nil {
				return nil, fmt.Errorf("key %q: cannot cast %q to %s: %w", k, v, rule.Type, err)
			}
			results = append(results, CastResult{Key: k, Original: v, Cast: cast, Type: rule.Type})
		}
	}
	return results, nil
}

func castValue(v string, t CastType) (interface{}, error) {
	switch t {
	case CastString:
		return v, nil
	case CastInt:
		return strconv.Atoi(strings.TrimSpace(v))
	case CastBool:
		return strconv.ParseBool(strings.TrimSpace(v))
	case CastFloat:
		return strconv.ParseFloat(strings.TrimSpace(v), 64)
	default:
		return nil, fmt.Errorf("unknown cast type: %s", t)
	}
}
