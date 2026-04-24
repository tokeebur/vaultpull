package vault

import (
	"fmt"
	"os"
)

// PromoteFiles is a convenience wrapper that reads src and dst env files,
// applies PromoteSecrets, and writes the result to dstFile.
// If dstFile does not exist it is created.
func PromoteFiles(srcFile, dstFile string, opts PromoteOptions) (PromoteResult, error) {
	src, err := ParseEnvFile(srcFile)
	if err != nil {
		return PromoteResult{}, fmt.Errorf("promote: reading src %q: %w", srcFile, err)
	}

	dst, err := ParseEnvFile(dstFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return PromoteResult{}, fmt.Errorf("promote: reading dst %q: %w", dstFile, err)
		}
		dst = map[string]string{}
	}

	out, result, err := PromoteSecrets(src, dst, opts)
	if err != nil {
		return result, err
	}

	if opts.DryRun {
		return result, nil
	}

	if err := WriteEnvFile(dstFile, out); err != nil {
		return result, fmt.Errorf("promote: writing dst %q: %w", dstFile, err)
	}
	return result, nil
}
