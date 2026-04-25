package vault

import (
	"fmt"
	"strings"
)

// UppercaseReport holds the result of an uppercase normalization pass.
type UppercaseReport struct {
	Renamed  []string // keys that were renamed (value unchanged)
	Conflict []string // new uppercase key collided with an existing key
}

// NormalizeKeys returns a new map whose keys are uppercased.
// If an uppercase collision occurs and overwrite is false, the original key is
// kept and the conflict is recorded in the report. When overwrite is true the
// uppercased key wins.
func NormalizeKeys(secrets map[string]string, overwrite bool) (map[string]string, UppercaseReport) {
	out := make(map[string]string, len(secrets))
	var report UppercaseReport

	for k, v := range secrets {
		upper := strings.ToUpper(k)
		if upper == k {
			out[k] = v
			continue
		}
		// key needs renaming
		if _, exists := out[upper]; exists && !overwrite {
			// collision — keep original key
			out[k] = v
			report.Conflict = append(report.Conflict, k)
		} else {
			out[upper] = v
			report.Renamed = append(report.Renamed, k)
		}
	}
	return out, report
}

// FormatUppercaseReport returns a human-readable summary of the normalization.
func FormatUppercaseReport(r UppercaseReport) string {
	if len(r.Renamed) == 0 && len(r.Conflict) == 0 {
		return "all keys already uppercase — no changes"
	}
	var sb strings.Builder
	if len(r.Renamed) > 0 {
		sb.WriteString(fmt.Sprintf("renamed (%d):\n", len(r.Renamed)))
		for _, k := range r.Renamed {
			sb.WriteString(fmt.Sprintf("  %s → %s\n", k, strings.ToUpper(k)))
		}
	}
	if len(r.Conflict) > 0 {
		sb.WriteString(fmt.Sprintf("conflicts kept as-is (%d):\n", len(r.Conflict)))
		for _, k := range r.Conflict {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
