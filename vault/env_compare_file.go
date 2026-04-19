package vault

import (
	"fmt"
	"os"
	"strings"
)

// CompareFiles compares two .env files on disk and returns a CompareResult.
func CompareFiles(pathA, pathB string) (CompareResult, error) {
	a, err := ParseEnvFile(pathA)
	if err != nil && !os.IsNotExist(err) {
		return CompareResult{}, fmt.Errorf("reading %s: %w", pathA, err)
	}
	b, err := ParseEnvFile(pathB)
	if err != nil && !os.IsNotExist(err) {
		return CompareResult{}, fmt.Errorf("reading %s: %w", pathB, err)
	}
	return CompareSecrets(a, b), nil
}

// WriteCompareReport writes a compare result to a file.
func WriteCompareReport(path string, r CompareResult, labelA, labelB string) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Compare report: %s vs %s\n", labelA, labelB))
	sb.WriteString(FormatCompareResult(r, labelA, labelB))
	sb.WriteString(fmt.Sprintf("\n# Summary: %d only-%s, %d only-%s, %d differ, %d identical\n",
		len(r.OnlyInA), labelA, len(r.OnlyInB), labelB, len(r.Differ), len(r.Identical)))
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
