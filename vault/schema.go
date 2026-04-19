package vault

import (
	"fmt"
	"regexp"
	"strings"
)

// FieldSchema defines expected type and optional pattern for a secret key.
type FieldSchema struct {
	Key      string
	Type     string // "string", "int", "bool"
	Pattern  string
	Required bool
}

// SchemaFile holds a list of field schemas.
type SchemaFile struct {
	Fields []FieldSchema
}

// ValidateAgainstSchema checks secrets map against a SchemaFile.
// Returns a list of violations.
func ValidateAgainstSchema(secrets map[string]string, schema SchemaFile) []string {
	var violations []string

	required := map[string]FieldSchema{}
	for _, f := range schema.Fields {
		if f.Required {
			required[f.Key] = f
		}
	}

	for _, field := range schema.Fields {
		val, exists := secrets[field.Key]
		if !exists {
			if field.Required {
				violations = append(violations, fmt.Sprintf("missing required key: %s", field.Key))
			}
			continue
		}
		if field.Pattern != "" {
			matched, err := regexp.MatchString(field.Pattern, val)
			if err != nil || !matched {
				violations = append(violations, fmt.Sprintf("key %s does not match pattern %s", field.Key, field.Pattern))
			}
		}
		switch strings.ToLower(field.Type) {
		case "int":
			if !regexp.MustCompile(`^-?\d+$`).MatchString(val) {
				violations = append(violations, fmt.Sprintf("key %s expected int, got: %s", field.Key, val))
			}
		case "bool":
			l := strings.ToLower(val)
			if l != "true" && l != "false" {
				violations = append(violations, fmt.Sprintf("key %s expected bool, got: %s", field.Key, val))
			}
		}
	}
	return violations
}
