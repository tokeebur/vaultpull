package vault

import "fmt"

// ChainEntry represents a single step in a secret resolution chain.
type ChainEntry struct {
	Source string
	Key    string
	Value  string
}

// SecretChain holds an ordered list of sources to resolve secrets from.
type SecretChain struct {
	Entries []ChainEntry
}

// NewSecretChain initializes an empty SecretChain.
func NewSecretChain() *SecretChain {
	return &SecretChain{}
}

// AddSource appends secrets from a named source map to the chain.
func (sc *SecretChain) AddSource(name string, secrets map[string]string) {
	for k, v := range secrets {
		sc.Entries = append(sc.Entries, ChainEntry{Source: name, Key: k, Value: v})
	}
}

// Resolve returns the last value set for each key across all sources,
// along with which source it came from.
func (sc *SecretChain) Resolve() (map[string]string, map[string]string) {
	resolved := make(map[string]string)
	sources := make(map[string]string)
	for _, e := range sc.Entries {
		resolved[e.Key] = e.Value
		sources[e.Key] = e.Source
	}
	return resolved, sources
}

// Explain returns a human-readable string showing resolution order for a key.
func (sc *SecretChain) Explain(key string) string {
	out := fmt.Sprintf("Resolution chain for %q:\n", key)
	found := false
	for _, e := range sc.Entries {
		if e.Key == key {
			out += fmt.Sprintf("  [%s] %s = %s\n", e.Source, e.Key, e.Value)
			found = true
		}
	}
	if !found {
		out += "  (no entries found)\n"
	}
	return out
}
