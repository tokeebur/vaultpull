package vault

import (
	"fmt"
	"sort"
	"strings"
)

// FormatEncryptReport returns a human-readable summary of encryption results.
func FormatEncryptReport(results []EncryptResult) string {
	if len(results) == 0 {
		return "No keys encrypted."
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Encrypted %d key(s):\n", len(results)))
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("  + %s\n", r.Key))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatDecryptReport returns a summary of how many keys were decrypted.
func FormatDecryptReport(original, decrypted map[string]string) string {
	count := 0
	for k, v := range decrypted {
		if orig, ok := original[k]; ok && orig != v {
			count++
		}
	}
	if count == 0 {
		return "No keys decrypted."
	}
	return fmt.Sprintf("Decrypted %d key(s).", count)
}
