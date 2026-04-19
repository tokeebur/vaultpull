package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AliasMap maps alias names to secret keys.
type AliasMap map[string]string

// LoadAliasFile reads an alias file where each line is: alias=key
func LoadAliasFile(path string) (AliasMap, error) {
	am := make(AliasMap)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return am, nil
	}
	if err != nil {
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
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid alias line: %q", line)
		}
		am[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return am, scanner.Err()
}

// SaveAliasFile writes an AliasMap to a file.
func SaveAliasFile(path string, am AliasMap) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	for alias, key := range am {
		if _, err := fmt.Fprintf(f, "%s=%s\n", alias, key); err != nil {
			return err
		}
	}
	return nil
}

// ApplyAliases returns a new map with aliased keys added alongside originals.
// If an alias target key exists in secrets, the alias name is added as an extra entry.
func ApplyAliases(secrets map[string]string, am AliasMap) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for alias, key := range am {
		if val, ok := secrets[key]; ok {
			out[alias] = val
		}
	}
	return out
}
