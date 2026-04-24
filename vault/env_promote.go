package vault

import (
	"fmt"
	"sort"
	"strings"
)

// PromoteOptions controls how secrets are promoted between environments.
type PromoteOptions struct {
	Overwrite bool
	DryRun    bool
	Keys      []string // if non-empty, only promote these keys
}

// PromoteResult describes the outcome of a promotion operation.
type PromoteResult struct {
	Promoted  []string
	Skipped   []string
	Overwrote []string
}

// PromoteSecrets copies secrets from src into dst according to opts.
// Returns a PromoteResult describing what happened.
func PromoteSecrets(src, dst map[string]string, opts PromoteOptions) (map[string]string, PromoteResult, error) {
	result := PromoteResult{}
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			return nil, result, fmt.Errorf("promote: key %q not found in source", k)
		}
		if _, exists := out[k]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if _, exists := out[k]; exists && opts.Overwrite {
			result.Overwrote = append(result.Overwrote, k)
		} else {
			result.Promoted = append(result.Promoted, k)
		}
		if !opts.DryRun {
			out[k] = v
		}
	}
	return out, result, nil
}

// FormatPromoteReport returns a human-readable summary of a PromoteResult.
func FormatPromoteReport(r PromoteResult) string {
	var sb strings.Builder
	if len(r.Promoted) > 0 {
		sb.WriteString(fmt.Sprintf("promoted (%d): %s\n", len(r.Promoted), strings.Join(r.Promoted, ", ")))
	}
	if len(r.Overwrote) > 0 {
		sb.WriteString(fmt.Sprintf("overwrote (%d): %s\n", len(r.Overwrote), strings.Join(r.Overwrote, ", ")))
	}
	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("skipped (%d): %s\n", len(r.Skipped), strings.Join(r.Skipped, ", ")))
	}
	if sb.Len() == 0 {
		return "nothing to promote\n"
	}
	return sb.String()
}
