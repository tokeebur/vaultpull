package vault

import (
	"fmt"
	"io"
	"strings"
)

// ListGroups writes a human-readable summary of all groups to w.
func ListGroups(gm GroupMap, w io.Writer) {
	names := GroupNames(gm)
	if len(names) == 0 {
		fmt.Fprintln(w, "(no groups defined)")
		return
	}
	for _, name := range names {
		keys := gm[name]
		fmt.Fprintf(w, "%-20s %s\n", name, strings.Join(keys, ", "))
	}
}

// GroupsForKey returns all group names that contain the given key.
func GroupsForKey(gm GroupMap, key string) []string {
	var result []string
	for _, name := range GroupNames(gm) {
		for _, k := range gm[name] {
			if k == key {
				result = append(result, name)
				break
			}
		}
	}
	return result
}
