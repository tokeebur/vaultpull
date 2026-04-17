package vault

import (
	"fmt"
	"strings"
)

// NamespaceConfig holds Vault namespace settings.
type NamespaceConfig struct {
	Namespace string
	MountPath string
	SecretPath string
}

// ParseNamespacedPath splits a full path like "ns/mount/secret" into components.
func ParseNamespacedPath(fullPath string) (*NamespaceConfig, error) {
	parts := strings.SplitN(fullPath, "/", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid namespaced path %q: expected at least namespace/mount/secret", fullPath)
	}
	cfg := &NamespaceConfig{
		Namespace: parts[0],
		MountPath: parts[1],
	}
	if len(parts) == 3 {
		cfg.SecretPath = parts[2]
	}
	return cfg, nil
}

// BuildVaultPath constructs the KV v2 API path for a namespaced secret.
func BuildVaultPath(ns *NamespaceConfig) string {
	if ns.SecretPath == "" {
		return fmt.Sprintf("%s/data/%s", ns.MountPath, ns.MountPath)
	}
	return fmt.Sprintf("%s/data/%s", ns.MountPath, ns.SecretPath)
}

// ApplyNamespaceHeader returns a copy of headers map with X-Vault-Namespace set.
func ApplyNamespaceHeader(headers map[string]string, namespace string) map[string]string {
	result := make(map[string]string, len(headers)+1)
	for k, v := range headers {
		result[k] = v
	}
	if namespace != "" {
		result["X-Vault-Namespace"] = namespace
	}
	return result
}
