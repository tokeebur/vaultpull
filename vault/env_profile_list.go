package vault

import (
	"fmt"
	"io"
	"sort"
)

// ListProfiles writes a summary of all profiles and their keys to w.
func ListProfiles(profiles map[string]Profile, w io.Writer) {
	names := make([]string, 0, len(profiles))
	for n := range profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		p := profiles[n]
		keys := make([]string, 0, len(p.Values))
		for k := range p.Values {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Fprintf(w, "[%s] (%d keys)\n", n, len(keys))
		for _, k := range keys {
			fmt.Fprintf(w, "  %s=%s\n", k, p.Values[k])
		}
	}
}

// ProfileNames returns sorted profile names from a profiles map.
func ProfileNames(profiles map[string]Profile) []string {
	names := make([]string, 0, len(profiles))
	for n := range profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
