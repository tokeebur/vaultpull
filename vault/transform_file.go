package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadTransformFile loads transform rules from a file.
// Each non-comment line: PATTERN=TRANSFORM
func LoadTransformFile(path string) ([]TransformRule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open transform file: %w", err)
	}
	defer f.Close()

	var rules []TransformRule
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid transform rule: %q", line)
		}
		rules = append(rules, TransformRule{
			KeyPattern: strings.TrimSpace(parts[0]),
			Transform:  strings.TrimSpace(parts[1]),
		})
	}
	return rules, scanner.Err()
}

// SaveTransformFile saves transform rules to a file.
func SaveTransformFile(path string, rules []TransformRule) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create transform file: %w", err)
	}
	defer f.Close()

	fw := bufio.NewWriter(f)
	fmt.Fprintln(fw, "# vaultpull transform rules: PATTERN=TRANSFORM")
	for _, r := range rules {
		fmt.Fprintf(fw, "%s=%s\n", r.KeyPattern, r.Transform)
	}
	return fw.Flush()
}
