package vault

import (
	"fmt"
	"os"
	"strings"
)

// MergeStrategy defines how conflicts are resolved during merge.
type MergeStrategy int

const (
	// MergeStrategyKeepLocal keeps existing local values on conflict.
	MergeStrategyKeepLocal MergeStrategy = iota
	// MergeStrategyKeepRemote overwrites local values with remote on conflict.
	MergeStrategyKeepRemote
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Added    []string
	Updated  []string
	Skipped  []string
	Final    map[string]string
}

// MergeSecrets merges remote secrets into the existing env file using the given strategy.
// It returns a MergeResult describing what changed and writes the merged result to outputPath.
func MergeSecrets(outputPath string, remote map[string]string, strategy MergeStrategy) (*MergeResult, error) {
	local, err := ParseEnvFile(outputPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading existing env file: %w", err)
	}
	if local == nil {
		local = make(map[string]string)
	}

	result := &MergeResult{
		Final: make(map[string]string),
	}

	// Copy all local keys first.
	for k, v := range local {
		result.Final[k] = v
	}

	for k, remoteVal := range remote {
		if localVal, exists := local[k]; exists {
			if localVal == remoteVal {
				result.Final[k] = localVal
				continue
			}
			switch strategy {
			case MergeStrategyKeepLocal:
				result.Skipped = append(result.Skipped, k)
				result.Final[k] = localVal
			case MergeStrategyKeepRemote:
				result.Updated = append(result.Updated, k)
				result.Final[k] = remoteVal
			}
		} else {
			result.Added = append(result.Added, k)
			result.Final[k] = remoteVal
		}
	}

	var sb strings.Builder
	for _, k := range sortedKeys(result.Final) {
		fmt.Fprintf(&sb, "%s=%s\n", k, result.Final[k])
	}

	if err := os.WriteFile(outputPath, []byte(sb.String()), 0600); err != nil {
		return nil, fmt.Errorf("writing merged env file: %w", err)
	}

	return result, nil
}
