package vault

import (
	"fmt"
	"strings"
)

// CoerceRule defines how a secret key should be coerced in output.
// Action can be: "upper", "lower", "title", "trim", "quote", "unquote".
type CoerceRule struct {
	KeyPattern string
	Action     string
}

// CoerceResult records what happened to a single key during coercion.
type CoerceResult struct {
	Key      string
	Original string
	Coerced  string
	Changed  bool
}

// CoerceReport holds the full set of results from a CoerceSecrets call.
type CoerceReport struct {
	Results []CoerceResult
}

// CoerceSecrets applies a list of coercion rules to the provided secrets map
// and returns a new map with transformed values alongside a report of changes.
// Rules are applied in order; the first matching rule wins.
func CoerceSecrets(secrets map[string]string, rules []CoerceRule) (map[string]string, CoerceReport, error) {
	out := make(map[string]string, len(secrets))
	var report CoerceReport

	for k, v := range secrets {
		matched := false
		for _, rule := range rules {
			if !matchesPattern(rule.KeyPattern, k) {
				continue
			}
			coerced, err := applyCoerceAction(v, rule.Action)
			if err != nil {
				return nil, CoerceReport{}, fmt.Errorf("coerce key %q with action %q: %w", k, rule.Action, err)
			}
			out[k] = coerced
			report.Results = append(report.Results, CoerceResult{
				Key:      k,
				Original: v,
				Coerced:  coerced,
				Changed:  coerced != v,
			})
			matched = true
			break
		}
		if !matched {
			out[k] = v
		}
	}

	return out, report, nil
}

// applyCoerceAction transforms a single value according to the named action.
func applyCoerceAction(value, action string) (string, error) {
	switch strings.ToLower(action) {
	case "upper":
		return strings.ToUpper(value), nil
	case "lower":
		return strings.ToLower(value), nil
	case "title":
		return strings.Title(strings.ToLower(value)), nil //nolint:staticcheck
	case "trim":
		return strings.TrimSpace(value), nil
	case "quote":
		// Wrap in double-quotes if not already quoted.
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			return value, nil
		}
		return `"` + value + `"`, nil
	case "unquote":
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			return value[1 : len(value)-1], nil
		}
		return value, nil
	default:
		return "", fmt.Errorf("unknown coerce action %q", action)
	}
}

// FormatCoerceReport returns a human-readable summary of coercion changes.
func FormatCoerceReport(report CoerceReport) string {
	if len(report.Results) == 0 {
		return "no coercions applied\n"
	}
	var sb strings.Builder
	changed := 0
	for _, r := range report.Results {
		if r.Changed {
			fmt.Fprintf(&sb, "  ~ %s: %q -> %q\n", r.Key, r.Original, r.Coerced)
			changed++
		}
	}
	if changed == 0 {
		return "coercions applied — no values changed\n"
	}
	return sb.String()
}
