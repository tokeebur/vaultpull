package vault

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// FileDiffResult holds the result of comparing two .env files on disk.
type FileDiffResult struct {
	OnlyInA  map[string]string // keys present only in file A
	OnlyInB  map[string]string // keys present only in file B
	Differ   map[string][2]string // keys present in both but with different values [a, b]
	Identical map[string]string // keys with identical values in both files
}

// DiffFiles reads two .env files and returns a structured diff of their contents.
// Missing files are treated as empty sets of secrets.
func DiffFiles(pathA, pathB string) (FileDiffResult, error) {
	secretsA, err := ParseEnvFile(pathA)
	if err != nil && !os.IsNotExist(err) {
		return FileDiffResult{}, fmt.Errorf("reading %s: %w", pathA, err)
	}
	secretsB, err := ParseEnvFile(pathB)
	if err != nil && !os.IsNotExist(err) {
		return FileDiffResult{}, fmt.Errorf("reading %s: %w", pathB, err)
	}

	result := FileDiffResult{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Differ:    make(map[string][2]string),
		Identical: make(map[string]string),
	}

	// Check all keys in A
	for k, va := range secretsA {
		if vb, ok := secretsB[k]; ok {
			if va == vb {
				result.Identical[k] = va
			} else {
				result.Differ[k] = [2]string{va, vb}
			}
		} else {
			result.OnlyInA[k] = va
		}
	}

	// Keys only in B
	for k, vb := range secretsB {
		if _, ok := secretsA[k]; !ok {
			result.OnlyInB[k] = vb
		}
	}

	return result, nil
}

// FormatFileDiff formats a FileDiffResult into a human-readable string.
// Values are masked by default to avoid leaking secrets in output.
func FormatFileDiff(result FileDiffResult, showValues bool) string {
	var sb strings.Builder

	maskOrShow := func(v string) string {
		if showValues {
			return v
		}
		return "***"
	}

	keys := func(m map[string]string) []string {
		out := make([]string, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}

	if len(result.OnlyInA) > 0 {
		sb.WriteString("Only in A:\n")
		for _, k := range keys(result.OnlyInA) {
			fmt.Fprintf(&sb, "  - %s=%s\n", k, maskOrShow(result.OnlyInA[k]))
		}
	}

	if len(result.OnlyInB) > 0 {
		sb.WriteString("Only in B:\n")
		for _, k := range keys(result.OnlyInB) {
			fmt.Fprintf(&sb, "  + %s=%s\n", k, maskOrShow(result.OnlyInB[k]))
		}
	}

	if len(result.Differ) > 0 {
		sb.WriteString("Changed:\n")
		dkeys := make([]string, 0, len(result.Differ))
		for k := range result.Differ {
			dkeys = append(dkeys, k)
		}
		sort.Strings(dkeys)
		for _, k := range dkeys {
			pair := result.Differ[k]
			fmt.Fprintf(&sb, "  ~ %s: %s -> %s\n", k, maskOrShow(pair[0]), maskOrShow(pair[1]))
		}
	}

	if sb.Len() == 0 {
		sb.WriteString("No differences found.\n")
	}

	return sb.String()
}

// WriteFileDiffReport writes a formatted diff report to the given path.
func WriteFileDiffReport(path string, result FileDiffResult, showValues bool) error {
	report := FormatFileDiff(result, showValues)
	return os.WriteFile(path, []byte(report), 0o644)
}
