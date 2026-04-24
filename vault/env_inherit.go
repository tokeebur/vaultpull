package vault

import (
	"fmt"
	"sort"
	"strings"
)

// InheritRule defines a key inheritance mapping from a parent environment.
type InheritRule struct {
	Key      string
	ParentKey string // if empty, same as Key
	Override bool   // if true, parent value overwrites child value
}

// InheritResult holds the outcome of a single key inheritance.
type InheritResult struct {
	Key      string
	Value    string
	Source   string // "child", "parent", or "default"
	Overridden bool
}

// InheritSecrets merges secrets from a parent map into a child map according
// to the provided rules. If a rule has Override=true, the parent value always
// wins. Otherwise, the child value is preserved when present.
func InheritSecrets(child, parent map[string]string, rules []InheritRule) (map[string]string, []InheritResult, error) {
	if child == nil {
		child = map[string]string{}
	}
	result := make(map[string]string, len(child))
	for k, v := range child {
		result[k] = v
	}

	var report []InheritResult

	for _, rule := range rules {
		if rule.Key == "" {
			return nil, nil, fmt.Errorf("inherit rule has empty Key")
		}
		parentKey := rule.ParentKey
		if parentKey == "" {
			parentKey = rule.Key
		}
		parentVal, parentHas := parent[parentKey]
		childVal, childHas := result[rule.Key]

		switch {
		case rule.Override && parentHas:
			result[rule.Key] = parentVal
			report = append(report, InheritResult{Key: rule.Key, Value: parentVal, Source: "parent", Overridden: childHas})
		case !childHas && parentHas:
			result[rule.Key] = parentVal
			report = append(report, InheritResult{Key: rule.Key, Value: parentVal, Source: "parent", Overridden: false})
		case childHas:
			report = append(report, InheritResult{Key: rule.Key, Value: childVal, Source: "child", Overridden: false})
		default:
			report = append(report, InheritResult{Key: rule.Key, Value: "", Source: "missing", Overridden: false})
		}
	}

	return result, report, nil
}

// FormatInheritReport returns a human-readable summary of inheritance results.
func FormatInheritReport(report []InheritResult) string {
	if len(report) == 0 {
		return "no inheritance rules applied\n"
	}
	sort.Slice(report, func(i, j int) bool {
		return report[i].Key < report[j].Key
	})
	var sb strings.Builder
	for _, r := range report {
		override := ""
		if r.Overridden {
			override = " (overridden)"
		}
		fmt.Fprintf(&sb, "%-30s source=%-8s%s\n", r.Key, r.Source, override)
	}
	return sb.String()
}
