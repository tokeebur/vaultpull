package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPolicyAllows_Default(t *testing.T) {
	p := NewPolicy("test", nil)
	ok, err := p.Allows("ANY_KEY")
	if err != nil || !ok {
		t.Fatalf("expected allow by default, got ok=%v err=%v", ok, err)
	}
}

func TestPolicyAllows_DenyPattern(t *testing.T) {
	p := NewPolicy("test", []PolicyRule{{Pattern: "^SECRET_", Allow: false}})
	ok, err := p.Allows("SECRET_TOKEN")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected deny for SECRET_TOKEN")
	}
	ok2, _ := p.Allows("DB_PASS")
	if !ok2 {
		t.Fatal("expected allow for DB_PASS")
	}
}

func TestPolicyAllows_LastMatchWins(t *testing.T) {
	p := NewPolicy("test", []PolicyRule{
		{Pattern: "^DB_", Allow: false},
		{Pattern: "^DB_PASS$", Allow: true},
	})
	ok, _ := p.Allows("DB_PASS")
	if !ok {
		t.Fatal("expected allow: last rule should win")
	}
	ok2, _ := p.Allows("DB_HOST")
	if ok2 {
		t.Fatal("expected deny for DB_HOST")
	}
}

func TestFilterByPolicy(t *testing.T) {
	p := NewPolicy("test", []PolicyRule{{Pattern: "^INTERNAL_", Allow: false}})
	secrets := map[string]string{"INTERNAL_KEY": "x", "PUBLIC_KEY": "y"}
	result, denied, err := FilterByPolicy(secrets, p)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["PUBLIC_KEY"]; !ok {
		t.Fatal("expected PUBLIC_KEY in result")
	}
	if len(denied) != 1 || denied[0] != "INTERNAL_KEY" {
		t.Fatalf("expected INTERNAL_KEY denied, got %v", denied)
	}
}

func TestSaveAndLoadPolicyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.policy")
	p := NewPolicy("test", []PolicyRule{
		{Pattern: "^DB_", Allow: true},
		{Pattern: "^SECRET_", Allow: false},
	})
	if err := SavePolicyFile(path, p); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadPolicyFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(loaded.Rules))
	}
}

func TestLoadPolicyFile_NotExist(t *testing.T) {
	_, err := LoadPolicyFile("/nonexistent/policy")
	if !os.IsNotExist(err) {
		// wrapped error
		if err == nil {
			t.Fatal("expected error")
		}
	}
}

func TestParsePolicyLine_Invalid(t *testing.T) {
	_, err := ParsePolicyLine("NOPREFIX")
	if err == nil {
		t.Fatal("expected error for missing prefix")
	}
}
