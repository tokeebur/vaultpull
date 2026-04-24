package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadScopeFile parses a scope definition file.
// Format:
//   [scope-name]
//   include=PATTERN
//   exclude=PATTERN
func LoadScopeFile(path string) ([]Scope, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var scopes []Scope
	var current *Scope
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if current != nil {
				scopes = append(scopes, *current)
			}
			current = &Scope{Name: line[1 : len(line)-1]}
			continue
		}
		if current == nil {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "include":
			current.Include = append(current.Include, val)
		case "exclude":
			current.Exclude = append(current.Exclude, val)
		}
	}
	if current != nil {
		scopes = append(scopes, *current)
	}
	return scopes, scanner.Err()
}

// SaveScopeFile writes scopes to a file in the standard format.
func SaveScopeFile(path string, scopes []Scope) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, s := range scopes {
		fmt.Fprintf(w, "[%s]\n", s.Name)
		for _, inc := range s.Include {
			fmt.Fprintf(w, "include=%s\n", inc)
		}
		for _, exc := range s.Exclude {
			fmt.Fprintf(w, "exclude=%s\n", exc)
		}
		fmt.Fprintln(w)
	}
	return w.Flush()
}
