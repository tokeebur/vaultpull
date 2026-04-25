package vault

import (
	"fmt"
	"strings"
)

// TrimRule defines how a secret key's value should be trimmed.
type TrimRule struct {
	Key    string // glob pattern or exact key
	Mode   string // "space", "prefix", "suffix", "both" (default: "space")
	Cutset string // used with prefix/suffix/both modes
}

// TrimReport records what was changed during trimming.
type TrimReport struct {
	Key    string
	Before string
	After  string
}

// TrimSecrets applies trim rules to a copy of secrets and returns the
// modified map along with a report of all changes made.
func TrimSecrets(secrets map[string]string, rules []TrimRule) (map[string]string, []TrimReport, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var report []TrimReport

	for _, rule := range rules {
		for k, v := range out {
			if !matchesPattern(rule.Key, k) {
				continue
			}
			var trimmed string
			switch rule.Mode {
			case "", "space":
				trimmed = strings.TrimSpace(v)
			case "prefix":
				trimmed = strings.TrimPrefix(v, rule.Cutset)
			case "suffix":
				trimmed = strings.TrimSuffix(v, rule.Cutset)
			case "both":
				trimmed = strings.Trim(v, rule.Cutset)
			default:
				return nil, nil, fmt.Errorf("unknown trim mode %q for key %q", rule.Mode, k)
			}
			if trimmed != v {
				report = append(report, TrimReport{Key: k, Before: v, After: trimmed})
				out[k] = trimmed
			}
		}
	}

	return out, report, nil
}

// FormatTrimReport returns a human-readable summary of trim changes.
func FormatTrimReport(report []TrimReport) string {
	if len(report) == 0 {
		return "no values trimmed\n"
	}
	var sb strings.Builder
	for _, r := range report {
		fmt.Fprintf(&sb, "trimmed %s: %q -> %q\n", r.Key, r.Before, r.After)
	}
	return sb.String()
}
