package vault

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested key segments are flattened.
type FlattenOptions struct {
	Separator   string // separator used to join nested segments (default "_")
	Uppercase   bool   // convert resulting keys to uppercase
	StripPrefix string // strip this prefix from each key before flattening
}

// FlattenSecrets takes a map whose keys may contain a separator (e.g. "." or "/")
// and re-joins the segments using opts.Separator, optionally uppercasing the result.
// The source separator is auto-detected as the first of ".", "/", ":" found in any key.
func FlattenSecrets(secrets map[string]string, opts FlattenOptions) (map[string]string, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	srcSep := detectSeparator(secrets)
	out := make(map[string]string, len(secrets))

	for k, v := range secrets {
		key := k
		if opts.StripPrefix != "" {
			key = strings.TrimPrefix(key, opts.StripPrefix)
		}
		if srcSep != "" && srcSep != opts.Separator {
			parts := strings.Split(key, srcSep)
			key = strings.Join(parts, opts.Separator)
		}
		if opts.Uppercase {
			key = strings.ToUpper(key)
		}
		if _, exists := out[key]; exists {
			return nil, fmt.Errorf("flatten: key collision on %q", key)
		}
		out[key] = v
	}
	return out, nil
}

// FormatFlattenReport returns a human-readable summary of the flatten operation.
func FormatFlattenReport(original, flattened map[string]string) string {
	var sb strings.Builder
	origKeys := sortedMapKeys(original)
	sb.WriteString(fmt.Sprintf("Flattened %d keys:\n", len(flattened)))
	for _, ok := range origKeys {
		nk := findNewKey(ok, original, flattened)
		if nk != ok {
			sb.WriteString(fmt.Sprintf("  %s -> %s\n", ok, nk))
		}
	}
	return sb.String()
}

func detectSeparator(secrets map[string]string) string {
	for k := range secrets {
		for _, sep := range []string{".", "/", ":"} {
			if strings.Contains(k, sep) {
				return sep
			}
		}
	}
	return ""
}

func sortedMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func findNewKey(origKey string, original, flattened map[string]string) string {
	origVal := original[origKey]
	for nk, nv := range flattened {
		if nv == origVal {
			return nk
		}
	}
	return origKey
}
