package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Profile represents a named set of secret overrides or mappings.
type Profile struct {
	Name   string
	Values map[string]string
}

// LoadProfileFile loads a profile file where sections are denoted by [profile-name].
func LoadProfileFile(path string) (map[string]Profile, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]Profile{}, nil
		}
		return nil, err
	}
	defer f.Close()

	profiles := map[string]Profile{}
	var current string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			current = line[1 : len(line)-1]
			if _, ok := profiles[current]; !ok {
				profiles[current] = Profile{Name: current, Values: map[string]string{}}
			}
			continue
		}
		if current == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		profiles[current].Values[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return profiles, scanner.Err()
}

// SaveProfileFile writes profiles to a file.
func SaveProfileFile(path string, profiles map[string]Profile) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	for name, p := range profiles {
		fmt.Fprintf(f, "[%s]\n", name)
		for k, v := range p.Values {
			fmt.Fprintf(f, "%s=%s\n", k, v)
		}
		fmt.Fprintln(f)
	}
	return nil
}

// ApplyProfile overlays profile values onto a secrets map.
func ApplyProfile(secrets map[string]string, p Profile) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}
	for k, v := range p.Values {
		result[k] = v
	}
	return result
}
