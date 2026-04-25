package vault

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Annotation holds metadata attached to a secret key.
type Annotation struct {
	Key   string
	Notes []string
}

// AnnotationMap maps secret keys to their annotations.
type AnnotationMap map[string][]string

// LoadAnnotationFile reads annotations from a file.
// Each non-comment line has the form: KEY note text here
func LoadAnnotationFile(path string) (AnnotationMap, error) {
	am := make(AnnotationMap)
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
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		note := strings.TrimSpace(parts[1])
		am[key] = append(am[key], note)
	}
	return am, scanner.Err()
}

// SaveAnnotationFile writes annotations to a file.
func SaveAnnotationFile(path string, am AnnotationMap) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	keys := make([]string, 0, len(am))
	for k := range am {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		for _, note := range am[k] {
			fmt.Fprintf(f, "%s %s\n", k, note)
		}
	}
	return nil
}

// AddAnnotation appends a note to the given key in the map.
func AddAnnotation(am AnnotationMap, key, note string) AnnotationMap {
	am[key] = append(am[key], note)
	return am
}

// FormatAnnotations returns a human-readable string of annotations for secrets.
func FormatAnnotations(secrets map[string]string, am AnnotationMap) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		notes, ok := am[k]
		if !ok || len(notes) == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s:\n", k))
		for _, n := range notes {
			sb.WriteString(fmt.Sprintf("  - %s\n", n))
		}
	}
	return sb.String()
}
