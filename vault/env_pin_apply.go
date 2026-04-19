package vault

import (
	"fmt"
	"strings"
)

// PinReport describes what was overridden by pins.
type PinReport struct {
	Key      string
	OldValue string
	NewValue string
	Pinned   bool
}

// ApplyPinsWithReport applies pins and returns a report of changes.
func ApplyPinsWithReport(secrets map[string]string, pf PinFile) (map[string]string, []PinReport) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	var reports []PinReport
	for key, entry := range pf {
		old := out[key]
		out[key] = entry.Value
		reports = append(reports, PinReport{
			Key:      key,
			OldValue: old,
			NewValue: entry.Value,
			Pinned:   true,
		})
	}
	return out, reports
}

// FormatPinReport returns a human-readable summary of pin overrides.
func FormatPinReport(reports []PinReport) string {
	if len(reports) == 0 {
		return "No pins applied."
	}
	var sb strings.Builder
	for _, r := range reports {
		if r.OldValue == "" {
			fmt.Fprintf(&sb, "[PIN] %s = %q (new key)\n", r.Key, r.NewValue)
		} else {
			fmt.Fprintf(&sb, "[PIN] %s: %q -> %q\n", r.Key, r.OldValue, r.NewValue)
		}
	}
	return sb.String()
}
