package vault

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR} and $VAR style references
var varPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// ExpandSecrets resolves variable references within secret values.
// References can point to other secrets or to OS environment variables.
// If a referenced key is not found and strict is true, an error is returned.
func ExpandSecrets(secrets map[string]string, strict bool) (map[string]string, error) {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		expanded, err := expandValue(v, secrets, strict)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		result[k] = expanded
	}
	return result, nil
}

// ListUnresolved returns keys whose values contain unresolvable references.
func ListUnresolved(secrets map[string]string) []string {
	var unresolved []string
	for k, v := range secrets {
		_, err := expandValue(v, secrets, true)
		if err != nil {
			unresolved = append(unresolved, k)
		}
	}
	return unresolved
}

func expandValue(val string, secrets map[string]string, strict bool) (string, error) {
	var expandErr error
	result := varPattern.ReplaceAllStringFunc(val, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := strings.TrimPrefix(strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}"), "$")
		if v, ok := secrets[name]; ok {
			return v
		}
		if v, ok := os.LookupEnv(name); ok {
			return v
		}
		if strict {
			expandErr = fmt.Errorf("unresolved reference: %s", name)
			return match
		}
		return match
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}
