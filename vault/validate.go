package vault

import (
	"fmt"
	"os"
	"regexp"
)

var validKeyPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// ValidationError holds a list of invalid keys found during validation.
type ValidationError struct {
	InvalidKeys []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid secret keys: %v", e.InvalidKeys)
}

// ValidateSecretKeys checks that all keys in secrets conform to env-var naming
// conventions (uppercase letters, digits, underscores; must not start with digit).
func ValidateSecretKeys(secrets map[string]interface{}) error {
	var invalid []string
	for k := range secrets {
		if !validKeyPattern.MatchString(k) {
			invalid = append(invalid, k)
		}
	}
	if len(invalid) > 0 {
		return &ValidationError{InvalidKeys: invalid}
	}
	return nil
}

// ValidateEnvFile parses an existing .env file and validates all key names.
func ValidateEnvFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}

	parsed, err := ParseEnvFile(path)
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	secrets := make(map[string]interface{}, len(parsed))
	for k, v := range parsed {
		secrets[k] = v
	}

	return ValidateSecretKeys(secrets)
}
