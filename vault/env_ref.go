package vault

import (
	"fmt"
	"sort"
	"strings"
)

// RefEntry describes a cross-reference between a secret key and an external source.
type RefEntry struct {
	Key    string
	Source string // e.g. "vault:secret/app#DB_PASS" or "env:DB_PASS"
	Note   string
}

// BuildRefIndex returns a map of key -> RefEntry for all keys that appear in refs.
func BuildRefIndex(refs []RefEntry) map[string]RefEntry {
	idx := make(map[string]RefEntry, len(refs))
	for _, r := range refs {
		idx[r.Key] = r
	}
	return idx
}

// ResolveRefs replaces secret values with the value pointed to by their ref source
// when the source is "literal:<value>". For real vault/env sources callers should
// pre-resolve before passing secrets in.
func ResolveRefs(secrets map[string]string, refs []RefEntry) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, r := range refs {
		if !strings.HasPrefix(r.Source, "literal:") {
			continue
		}
		val := strings.TrimPrefix(r.Source, "literal:")
		out[r.Key] = val
	}
	return out, nil
}

// FormatRefReport returns a human-readable summary of ref entries.
func FormatRefReport(refs []RefEntry) string {
	if len(refs) == 0 {
		return "no refs defined\n"
	}
	sorted := make([]RefEntry, len(refs))
	copy(sorted, refs)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })
	var sb strings.Builder
	for _, r := range sorted {
		line := fmt.Sprintf("%-24s -> %s", r.Key, r.Source)
		if r.Note != "" {
			line += fmt.Sprintf(" (%s)", r.Note)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}
