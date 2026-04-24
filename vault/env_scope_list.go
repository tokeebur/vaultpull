package vault

import (
	"fmt"
	"sort"
	"strings"
)

// ScopeNames returns the sorted list of scope names from a slice of Scope.
func ScopeNames(scopes []Scope) []string {
	names := make([]string, len(scopes))
	for i, s := range scopes {
		names[i] = s.Name
	}
	sort.Strings(names)
	return names
}

// FindScope returns the first scope matching name, or nil if not found.
func FindScope(scopes []Scope, name string) *Scope {
	for _, s := range scopes {
		if s.Name == name {
			copy := s
			return &copy
		}
	}
	return nil
}

// ListScopes returns a formatted string listing all scopes and their rules.
func ListScopes(scopes []Scope) string {
	if len(scopes) == 0 {
		return "no scopes defined"
	}
	sorted := make([]Scope, len(scopes))
	copy(sorted, scopes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	var sb strings.Builder
	for _, s := range sorted {
		fmt.Fprintf(&sb, "[%s]\n", s.Name)
		for _, inc := range s.Include {
			fmt.Fprintf(&sb, "  include: %s\n", inc)
		}
		for _, exc := range s.Exclude {
			fmt.Fprintf(&sb, "  exclude: %s\n", exc)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
