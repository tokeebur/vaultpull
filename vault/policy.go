package vault

import (
	"fmt"
	"regexp"
	"strings"
)

// PolicyRule defines an allow/deny rule for secret key access.
type PolicyRule struct {
	Pattern string
	Allow   bool
}

// Policy holds a named set of rules.
type Policy struct {
	Name  string
	Rules []PolicyRule
}

// NewPolicy creates a Policy with the given name and rules.
func NewPolicy(name string, rules []PolicyRule) *Policy {
	return &Policy{Name: name, Rules: rules}
}

// Allows returns true if the key is permitted by the policy.
// Rules are evaluated in order; last match wins. Default is allow.
func (p *Policy) Allows(key string) (bool, error) {
	allowed := true
	for _, rule := range p.Rules {
		matched, err := regexp.MatchString(rule.Pattern, key)
		if err != nil {
			return false, fmt.Errorf("invalid pattern %q: %w", rule.Pattern, err)
		}
		if matched {
			allowed = rule.Allow
		}
	}
	return allowed, nil
}

// FilterByPolicy returns only the secrets permitted by the policy.
func FilterByPolicy(secrets map[string]string, p *Policy) (map[string]string, []string, error) {
	result := make(map[string]string)
	var denied []string
	for k, v := range secrets {
		ok, err := p.Allows(k)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			result[k] = v
		} else {
			denied = append(denied, k)
		}
	}
	return result, denied, nil
}

// ParsePolicyLine parses a line like "+PATTERN" or "-PATTERN".
func ParsePolicyLine(line string) (PolicyRule, error) {
	line = strings.TrimSpace(line)
	if len(line) < 2 {
		return PolicyRule{}, fmt.Errorf("invalid policy line: %q", line)
	}
	switch line[0] {
	case '+':
		return PolicyRule{Pattern: line[1:], Allow: true}, nil
	case '-':
		return PolicyRule{Pattern: line[1:], Allow: false}, nil
	default:
		return PolicyRule{}, fmt.Errorf("policy line must start with + or -: %q", line)
	}
}
