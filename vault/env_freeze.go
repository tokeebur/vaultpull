package vault

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// FreezeEntry represents a frozen secret key with metadata.
type FreezeEntry struct {
	Key       string
	FrozenAt  time.Time
	FrozenBy  string
}

// FreezeFile maps keys to their freeze entries.
type FreezeFile map[string]FreezeEntry

// LoadFreezeFile reads a freeze file from disk.
// Each line: KEY=frozenAt|frozenBy
// Lines starting with '#' are skipped.
func LoadFreezeFile(path string) (FreezeFile, error) {
	ff := make(FreezeFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ff, nil
		}
		return nil, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		key := line[:eq]
		val := line[eq+1:]
		parts := strings.SplitN(val, "|", 2)
		if len(parts) != 2 {
			continue
		}
		t, err := time.Parse(time.RFC3339, parts[0])
		if err != nil {
			continue
		}
		ff[key] = FreezeEntry{Key: key, FrozenAt: t, FrozenBy: parts[1]}
	}
	return ff, nil
}

// SaveFreezeFile writes a freeze file to disk.
func SaveFreezeFile(path string, ff FreezeFile) error {
	keys := make([]string, 0, len(ff))
	for k := range ff {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString("# vaultpull freeze file\n")
	for _, k := range keys {
		e := ff[k]
		fmt.Fprintf(&sb, "%s=%s|%s\n", k, e.FrozenAt.UTC().Format(time.RFC3339), e.FrozenBy)
	}
	return os.WriteFile(path, []byte(sb.String()), 0600)
}

// FreezeKeys marks the given keys as frozen in the freeze file.
func FreezeKeys(ff FreezeFile, keys []string, frozenBy string) FreezeFile {
	out := make(FreezeFile, len(ff))
	for k, v := range ff {
		out[k] = v
	}
	now := time.Now().UTC()
	for _, k := range keys {
		out[k] = FreezeEntry{Key: k, FrozenAt: now, FrozenBy: frozenBy}
	}
	return out
}

// FilterFrozen removes frozen keys from secrets, returning the filtered map
// and a list of keys that were dropped.
func FilterFrozen(secrets map[string]string, ff FreezeFile) (map[string]string, []string) {
	out := make(map[string]string, len(secrets))
	var dropped []string
	for k, v := range secrets {
		if _, frozen := ff[k]; frozen {
			dropped = append(dropped, k)
			continue
		}
		out[k] = v
	}
	sort.Strings(dropped)
	return out, dropped
}
