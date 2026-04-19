package vault

import (
	"bytes"
	"strings"
	"testing"
)

func TestListProfiles_Output(t *testing.T) {
	profiles := map[string]Profile{
		"dev":  {Name: "dev", Values: map[string]string{"A": "1"}},
		"prod": {Name: "prod", Values: map[string]string{"A": "2", "B": "3"}},
	}
	var buf bytes.Buffer
	ListProfiles(profiles, &buf)
	out := buf.String()
	if !strings.Contains(out, "[dev]") {
		t.Errorf("expected [dev] in output")
	}
	if !strings.Contains(out, "[prod] (2 keys)") {
		t.Errorf("expected [prod] (2 keys) in output")
	}
	if !strings.Contains(out, "A=2") {
		t.Errorf("expected A=2 in output")
	}
}

func TestListProfiles_Empty(t *testing.T) {
	var buf bytes.Buffer
	ListProfiles(map[string]Profile{}, &buf)
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty profiles")
	}
}

func TestProfileNames_Sorted(t *testing.T) {
	profiles := map[string]Profile{
		"zebra": {Name: "zebra"},
		"alpha": {Name: "alpha"},
		"beta":  {Name: "beta"},
	}
	names := ProfileNames(profiles)
	if names[0] != "alpha" || names[1] != "beta" || names[2] != "zebra" {
		t.Errorf("expected sorted names, got %v", names)
	}
}

func TestListProfiles_SortedOutput(t *testing.T) {
	profiles := map[string]Profile{
		"z": {Name: "z", Values: map[string]string{"K": "v"}},
		"a": {Name: "a", Values: map[string]string{"K": "v"}},
	}
	var buf bytes.Buffer
	ListProfiles(profiles, &buf)
	lines := strings.Split(buf.String(), "\n")
	if !strings.HasPrefix(lines[0], "[a]") {
		t.Errorf("expected [a] first, got %s", lines[0])
	}
}
