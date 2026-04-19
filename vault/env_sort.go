package vault

import (
	"fmt"
	"sort"
	"strings"
)

// SortOrder defines how secrets should be sorted.
type SortOrder string

const (
	SortAlpha    SortOrder = "alpha"
	SortAlphaDesc SortOrder = "alpha-desc"
	SortByLength SortOrder = "length"
)

// SortedSecret holds a key-value pair for sorted output.
type SortedSecret struct {
	Key   string
	Value string
}

// SortSecrets returns secrets sorted by the given order.
func SortSecrets(secrets map[string]string, order SortOrder) ([]SortedSecret, error) {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}

	switch order {
	case SortAlpha, "":
		sort.Strings(keys)
	case SortAlphaDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	case SortByLength:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) == len(keys[j]) {
				return keys[i] < keys[j]
			}
			return len(keys[i]) < len(keys[j])
		})
	default:
		return nil, fmt.Errorf("unknown sort order: %q", order)
	}

	result := make([]SortedSecret, 0, len(keys))
	for _, k := range keys {
		result = append(result, SortedSecret{Key: k, Value: secrets[k]})
	}
	return result, nil
}

// FormatSorted formats sorted secrets as KEY=VALUE lines.
func FormatSorted(sorted []SortedSecret) string {
	var sb strings.Builder
	for _, s := range sorted {
		sb.WriteString(fmt.Sprintf("%s=%s\n", s.Key, s.Value))
	}
	return sb.String()
}
