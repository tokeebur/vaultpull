package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadTagMap reads a tag map file where each line is KEY=tag1,tag2.
func LoadTagMap(path string) (map[string][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]string{}, nil
		}
		return nil, fmt.Errorf("open tag map: %w", err)
	}
	defer f.Close()

	result := map[string][]string{}
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
		tags := strings.Split(parts[1], ",")
		for i, t := range tags {
			tags[i] = strings.TrimSpace(t)
		}
		result[strings.TrimSpace(parts[0])] = tags
	}
	return result, scanner.Err()
}

// SaveTagMap writes a tag map to a file.
func SaveTagMap(path string, tagMap map[string][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create tag map: %w", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for k, tags := range tagMap {
		fmt.Fprintf(w, "%s=%s\n", k, strings.Join(tags, ","))
	}
	return w.Flush()
}

// TagsForKey returns the tags associated with the given key, or nil if the key
// is not present in the tag map.
func TagsForKey(tagMap map[string][]string, key string) []string {
	tags, ok := tagMap[key]
	if !ok {
		return nil
	}
	return tags
}

// AddTag adds a tag to the given key in the tag map if it is not already present.
func AddTag(tagMap map[string][]string, key, tag string) {
	for _, t := range tagMap[key] {
		if t == tag {
			return
		}
	}
	tagMap[key] = append(tagMap[key], tag)
}
