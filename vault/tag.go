package vault

import (
	"fmt"
	"sort"
	"strings"
)

// TaggedSecret represents a secret with associated tags.
type TaggedSecret struct {
	Key   string
	Value string
	Tags  []string
}

// TagSecrets assigns tags to secrets based on a tag map (key -> tags).
func TagSecrets(secrets map[string]string, tagMap map[string][]string) []TaggedSecret {
	result := make([]TaggedSecret, 0, len(secrets))
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		tags := tagMap[k]
		if tags == nil {
			tags = []string{}
		}
		result = append(result, TaggedSecret{Key: k, Value: secrets[k], Tags: tags})
	}
	return result
}

// FilterByTag returns only secrets that have the given tag.
func FilterByTag(tagged []TaggedSecret, tag string) []TaggedSecret {
	var out []TaggedSecret
	for _, ts := range tagged {
		for _, t := range ts.Tags {
			if strings.EqualFold(t, tag) {
				out = append(out, ts)
				break
			}
		}
	}
	return out
}

// FormatTagged returns a human-readable string of tagged secrets.
func FormatTagged(tagged []TaggedSecret) string {
	var sb strings.Builder
	for _, ts := range tagged {
		sb.WriteString(fmt.Sprintf("%s [%s]\n", ts.Key, strings.Join(ts.Tags, ",")))
	}
	return sb.String()
}
