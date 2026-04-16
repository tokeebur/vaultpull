package vault

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// ExportFormat defines supported export formats.
type ExportFormat string

const (
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
	FormatYAML   ExportFormat = "yaml"
)

// ExportSecrets writes secrets to a file in the specified format.
func ExportSecrets(secrets map[string]string, path string, format ExportFormat) error {
	var sb strings.Builder

	keys := sortedKeys(secrets)

	switch format {
	case FormatDotenv:
		for _, k := range keys {
			fmt.Fprintf(&sb, "%s=%q\n", k, secrets[k])
		}
	case FormatJSON:
		sb.WriteString("{\n")
		for i, k := range keys {
			comma := ","
			if i == len(keys)-1 {
				comma = ""
			}
			fmt.Fprintf(&sb, "  %q: %q%s\n", k, secrets[k], comma)
		}
		sb.WriteString("}\n")
	case FormatYAML:
		for _, k := range keys {
			fmt.Fprintf(&sb, "%s: %q\n", k, secrets[k])
		}
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}

	return os.WriteFile(path, []byte(sb.String()), 0600)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
