package vault

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// GroupMap maps group names to lists of secret keys.
type GroupMap map[string][]string

// LoadGroupFile loads a group definition file.
// Format: groupname=KEY1,KEY2,KEY3
func LoadGroupFile(path string) (GroupMap, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return GroupMap{}, nil
		}
		return nil, err
	}
	defer f.Close()

	gm := GroupMap{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		group := strings.TrimSpace(parts[0])
		keys := strings.Split(strings.TrimSpace(parts[1]), ",")
		for i, k := range keys {
			keys[i] = strings.TrimSpace(k)
		}
		gm[group] = keys
	}
	return gm, scanner.Err()
}

// SaveGroupFile writes a GroupMap to a file.
func SaveGroupFile(path string, gm GroupMap) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	names := make([]string, 0, len(gm))
	for k := range gm {
		names = append(names, k)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Fprintf(f, "%s=%s\n", name, strings.Join(gm[name], ","))
	}
	return nil
}

// FilterByGroup returns only the secrets whose keys belong to the given group.
func FilterByGroup(secrets map[string]string, gm GroupMap, group string) (map[string]string, error) {
	keys, ok := gm[group]
	if !ok {
		return nil, fmt.Errorf("group %q not found", group)
	}
	result := make(map[string]string)
	for _, k := range keys {
		if v, exists := secrets[k]; exists {
			result[k] = v
		}
	}
	return result, nil
}

// GroupNames returns sorted group names.
func GroupNames(gm GroupMap) []string {
	names := make([]string, 0, len(gm))
	for k := range gm {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
