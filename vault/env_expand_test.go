package vault

import (
	"os"
	"testing"
)

func TestExpandSecrets_NoRefs(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := ExpandSecrets(secrets, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestExpandSecrets_ResolvesInternalRef(t *testing.T) {
	secrets := map[string]string{
		"HOST": "localhost",
		"URL":  "http://${HOST}:8080",
	}
	out, err := ExpandSecrets(secrets, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://localhost:8080" {
		t.Errorf("expected expanded URL, got %q", out["URL"])
	}
}

func TestExpandSecrets_ResolvesEnvVar(t *testing.T) {
	os.Setenv("MY_HOST", "envhost")
	defer os.Unsetenv("MY_HOST")
	secrets := map[string]string{"ADDR": "${MY_HOST}:9000"}
	out, err := ExpandSecrets(secrets, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ADDR"] != "envhost:9000" {
		t.Errorf("expected envhost:9000, got %q", out["ADDR"])
	}
}

func TestExpandSecrets_StrictMissingRef(t *testing.T) {
	secrets := map[string]string{"URL": "http://${MISSING}:8080"}
	_, err := ExpandSecrets(secrets, true)
	if err == nil {
		t.Fatal("expected error for missing ref in strict mode")
	}
}

func TestExpandSecrets_NonStrictMissingRef(t *testing.T) {
	secrets := map[string]string{"URL": "http://${MISSING}:8080"}
	out, err := ExpandSecrets(secrets, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://${MISSING}:8080" {
		t.Errorf("expected unchanged value, got %q", out["URL"])
	}
}

func TestListUnresolved(t *testing.T) {
	secrets := map[string]string{
		"GOOD": "static",
		"BAD":  "${GHOST}",
	}
	unresolved := ListUnresolved(secrets)
	if len(unresolved) != 1 || unresolved[0] != "BAD" {
		t.Errorf("expected [BAD], got %v", unresolved)
	}
}
