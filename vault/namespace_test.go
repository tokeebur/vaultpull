package vault

import (
	"testing"
)

func TestParseNamespacedPath_Valid(t *testing.T) {
	nc, err := ParseNamespacedPath("myns/secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nc.Namespace != "myns" {
		t.Errorf("expected namespace myns, got %s", nc.Namespace)
	}
	if nc.MountPath != "secret" {
		t.Errorf("expected mount secret, got %s", nc.MountPath)
	}
	if nc.SecretPath != "myapp" {
		t.Errorf("expected secret path myapp, got %s", nc.SecretPath)
	}
}

func TestParseNamespacedPath_TwoParts(t *testing.T) {
	nc, err := ParseNamespacedPath("myns/secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nc.SecretPath != "" {
		t.Errorf("expected empty secret path, got %s", nc.SecretPath)
	}
}

func TestParseNamespacedPath_Invalid(t *testing.T) {
	_, err := ParseNamespacedPath("onlyone")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestBuildVaultPath_WithSecret(t *testing.T) {
	nc := &NamespaceConfig{Namespace: "myns", MountPath: "secret", SecretPath: "myapp"}
	path := BuildVaultPath(nc)
	expected := "secret/data/myapp"
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestBuildVaultPath_NoSecret(t *testing.T) {
	nc := &NamespaceConfig{Namespace: "myns", MountPath: "secret", SecretPath: ""}
	path := BuildVaultPath(nc)
	expected := "secret/data/secret"
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestApplyNamespaceHeader(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer token"}
	result := ApplyNamespaceHeader(headers, "myns")
	if result["X-Vault-Namespace"] != "myns" {
		t.Errorf("expected X-Vault-Namespace myns, got %s", result["X-Vault-Namespace"])
	}
	if result["Authorization"] != "Bearer token" {
		t.Error("original headers should be preserved")
	}
	if len(headers) != 1 {
		t.Error("original map should not be mutated")
	}
}
