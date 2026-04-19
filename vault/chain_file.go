package vault

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ChainConfig describes an ordered list of named sources and their paths.
type ChainConfig struct {
	Sources []ChainSource
}

// ChainSource is a named secret source with a Vault path or local file path.
type ChainSource struct {
	Name string
	Path string
}

// LoadChainFile reads a chain config file with lines: name=path
func LoadChainFile(filename string) (*ChainConfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open chain file: %w", err)
	}
	defer f.Close()

	cc := &ChainConfig{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid chain line: %q", line)
		}
		cc.Sources = append(cc.Sources, ChainSource{
			Name: strings.TrimSpace(parts[0]),
			Path: strings.TrimSpace(parts[1]),
		})
	}
	return cc, scanner.Err()
}

// SaveChainFile writes a ChainConfig to a file.
func SaveChainFile(filename string, cc *ChainConfig) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create chain file: %w", err)
	}
	defer f.Close()
	for _, s := range cc.Sources {
		fmt.Fprintf(f, "%s=%s\n", s.Name, s.Path)
	}
	return nil
}
