package vault

import (
	"regexp"
	"strings"
)

// FilterMode controls how keys are matched.
type FilterMode string

const (
	FilterModePrefix FilterMode = "prefix"
	FilterModeSuffix FilterMode = "suffix"
	FilterModeRegex  FilterMode = "regex"
	FilterModeExact  FilterMode = "exact"
)

// FilterOptions configures secret filtering.
type FilterOptions struct {
	Mode    FilterMode
	Pattern string
	Invert  bool
}

// FilterSecrets returns a subset of secrets whose keys match the given options.
func FilterSecrets(secrets map[string]string, opts FilterOptions) (map[string]string, error) {
	result := make(map[string]string)

	var matchFn func(key string) (bool, error)

	switch opts.Mode {
	case FilterModePrefix:
		matchFn = func(key string) (bool, error) {
			return strings.HasPrefix(key, opts.Pattern), nil
		}
	case FilterModeSuffix:
		matchFn = func(key string) (bool, error) {
			return strings.HasSuffix(key, opts.Pattern), nil
		}
	case FilterModeRegex:
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
		matchFn = func(key string) (bool, error) {
			return re.MatchString(key), nil
		}
	default: // exact
		matchFn = func(key string) (bool, error) {
			return key == opts.Pattern, nil
		}
	}

	for k, v := range secrets {
		matched, err := matchFn(k)
		if err != nil {
			return nil, err
		}
		if opts.Invert {
			matched = !matched
		}
		if matched {
			result[k] = v
		}
	}
	return result, nil
}
