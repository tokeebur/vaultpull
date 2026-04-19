package vault

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a secret value.
type TransformFunc func(string) (string, error)

// TransformRule maps a key pattern to a named transform.
type TransformRule struct {
	KeyPattern string
	Transform  string
}

// BuiltinTransforms contains the available named transforms.
var BuiltinTransforms = map[string]TransformFunc{
	"upper":   func(v string) (string, error) { return strings.ToUpper(v), nil },
	"lower":   func(v string) (string, error) { return strings.ToLower(v), nil },
	"trim":    func(v string) (string, error) { return strings.TrimSpace(v), nil },
	"base64":  Base64Encode,
	"base64d": Base64Decode,
}

// ApplyTransforms applies matching transform rules to secrets.
func ApplyTransforms(secrets map[string]string, rules []TransformRule) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, rule := range rules {
		fn, ok := BuiltinTransforms[rule.Transform]
		if !ok {
			return nil, fmt.Errorf("unknown transform: %s", rule.Transform)
		}
		for k, v := range out {
			if matchesPattern(k, rule.KeyPattern) {
				result, err := fn(v)
				if err != nil {
					return nil, fmt.Errorf("transform %s failed on key %s: %w", rule.Transform, k, err)
				}
				out[k] = result
			}
		}
	}
	return out, nil
}

func matchesPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}
	return key == pattern
}
