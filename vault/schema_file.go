package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadSchemaFile parses a simple schema definition file.
// Format per line: KEY type required pattern
// Example: DB_PORT int true ^\d+$
func LoadSchemaFile(path string) (SchemaFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return SchemaFile{}, fmt.Errorf("open schema: %w", err)
	}
	defer f.Close()

	var schema SchemaFile
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		field := FieldSchema{
			Key:  parts[0],
			Type: parts[1],
		}
		if len(parts) >= 3 {
			field.Required = strings.EqualFold(parts[2], "true")
		}
		if len(parts) >= 4 {
			field.Pattern = parts[3]
		}
		schema.Fields = append(schema.Fields, field)
	}
	if err := scanner.Err(); err != nil {
		return SchemaFile{}, fmt.Errorf("scan schema: %w", err)
	}
	return schema, nil
}

// SaveSchemaFile writes a SchemaFile to disk.
func SaveSchemaFile(path string, schema SchemaFile) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create schema file: %w", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "# vaultpull schema file: KEY type required pattern")
	for _, field := range schema.Fields {
		req := "false"
		if field.Required {
			req = "true"
		}
		pat := field.Pattern
		if pat == "" {
			pat = "-"
		}
		fmt.Fprintf(w, "%s %s %s %s\n", field.Key, field.Type, req, pat)
	}
	return w.Flush()
}
