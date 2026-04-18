package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadPolicyFile reads a policy file and returns a Policy.
// Lines starting with '#' or empty lines are ignored.
// Each rule line must start with '+' (allow) or '-' (deny).
func LoadPolicyFile(path string) (*Policy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open policy file: %w", err)
	}
	defer f.Close()

	name := strings.TrimSuffix(path, ".policy")
	var rules []PolicyRule
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rule, err := ParsePolicyLine(line)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan policy file: %w", err)
	}
	return NewPolicy(name, rules), nil
}

// SavePolicyFile writes a Policy to a file.
func SavePolicyFile(path string, p *Policy) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create policy file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintf(w, "# Policy: %s\n", p.Name)
	for _, rule := range p.Rules {
		prefix := "+"
		if !rule.Allow {
			prefix = "-"
		}
		fmt.Fprintf(w, "%s%s\n", prefix, rule.Pattern)
	}
	return w.Flush()
}
