package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPolicyFile_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.policy")
	content := "# this is a comment\n\n+^API_\n-^INTERNAL_\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	p, err := LoadPolicyFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(p.Rules))
	}
	if p.Rules[0].Pattern != "^API_" || !p.Rules[0].Allow {
		t.Errorf("unexpected first rule: %+v", p.Rules[0])
	}
	if p.Rules[1].Pattern != "^INTERNAL_" || p.Rules[1].Allow {
		t.Errorf("unexpected second rule: %+v", p.Rules[1])
	}
}

func TestSavePolicyFile_Content(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.policy")
	p := NewPolicy("mypol", []PolicyRule{
		{Pattern: "^DB_", Allow: true},
		{Pattern: "^TMP_", Allow: false},
	})
	if err := SavePolicyFile(path, p); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !contains(s, "+^DB_") {
		t.Errorf("expected +^DB_ in output: %s", s)
	}
	if !contains(s, "-^TMP_") {
		t.Errorf("expected -^TMP_ in output: %s", s)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
