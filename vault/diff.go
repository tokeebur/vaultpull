package vault

import (
	"bufio"
	"os"
	"strings"
)

// EnvDiff represents the difference between existing and incoming secrets
type EnvDiff struct {
	Added   map[string]string
	Changed map[string]string
	Removed map[string]string
	Unchanged map[string]string
}

// ParseEnvFile reads an existing .env file into a key/value map
func ParseEnvFile(path string) (map[string]string, error) {
	result := make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result, scanner.Err()
}

// ComputeDiff compares existing env vars with incoming secrets
func ComputeDiff(existing, incoming map[string]string) EnvDiff {
	diff := EnvDiff{
		Added:     make(map[string]string),
		Changed:   make(map[string]string),
		Removed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, v := range incoming {
		if old, ok := existing[k]; !ok {
			diff.Added[k] = v
		} else if old != v {
			diff.Changed[k] = v
		} else {
			diff.Unchanged[k] = v
		}
	}

	for k, v := range existing {
		if _, ok := incoming[k]; !ok {
			diff.Removed[k] = v
		}
	}

	return diff
}
