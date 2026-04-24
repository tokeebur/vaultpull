package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckDeprecations_NoMatches(t *testing.T) {
	secrets := map[string]string{"NEW_KEY": "value"}
	rules := map[string]DeprecationEntry{
		"OLD_KEY": {Key: "OLD_KEY", Replacement: "NEW_KEY"},
	}
	r := CheckDeprecations(secrets, rules)
	if len(r.Deprecated) != 0 {
		t.Fatalf("expected 0 deprecated, got %d", len(r.Deprecated))
	}
}

func TestCheckDeprecations_FindsDeprecated(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "v", "FINE": "ok"}
	rules := map[string]DeprecationEntry{
		"OLD_KEY": {Key: "OLD_KEY", Replacement: "NEW_KEY", Reason: "renamed"},
	}
	r := CheckDeprecations(secrets, rules)
	if len(r.Deprecated) != 1 {
		t.Fatalf("expected 1, got %d", len(r.Deprecated))
	}
	if r.Deprecated[0].Key != "OLD_KEY" {
		t.Errorf("unexpected key %s", r.Deprecated[0].Key)
	}
}

func TestFormatDeprecationReport_Empty(t *testing.T) {
	out := FormatDeprecationReport(DeprecationReport{})
	if !strings.Contains(out, "No deprecated") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatDeprecationReport_WithEntries(t *testing.T) {
	r := DeprecationReport{
		Deprecated: []DeprecationEntry{
			{Key: "OLD_DB_PASS", Replacement: "DB_PASSWORD", Reason: "standardised"},
		},
	}
	out := FormatDeprecationReport(r)
	if !strings.Contains(out, "OLD_DB_PASS") {
		t.Errorf("missing key in output: %s", out)
	}
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("missing replacement in output: %s", out)
	}
	if !strings.Contains(out, "standardised") {
		t.Errorf("missing reason in output: %s", out)
	}
}

func TestSaveAndLoadDeprecationFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "deprecations.conf")
	rules := map[string]DeprecationEntry{
		"OLD_KEY": {Key: "OLD_KEY", Replacement: "NEW_KEY", Reason: "renamed"},
		"LEGACY":  {Key: "LEGACY"},
	}
	if err := SaveDeprecationFile(p, rules); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadDeprecationFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := loaded["OLD_KEY"]; !ok {
		t.Error("OLD_KEY missing after round-trip")
	}
	if loaded["OLD_KEY"].Replacement != "NEW_KEY" {
		t.Errorf("replacement mismatch: %s", loaded["OLD_KEY"].Replacement)
	}
}

func TestLoadDeprecationFile_NotExist(t *testing.T) {
	rules, err := LoadDeprecationFile("/nonexistent/path/dep.conf")
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty map, got %d entries", len(rules))
	}
}

func TestLoadDeprecationFile_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "dep.conf")
	content := "# this is a comment\nOLD_FOO=NEW_FOO:old\n"
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	rules, err := LoadDeprecationFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}
