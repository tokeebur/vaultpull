package vault

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewClient_InvalidAddr(t *testing.T) {
	_, err := NewClient("", "token")
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestNewClient_ValidConfig(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestFetchSecrets_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "test-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"data":{"API_KEY":"abc123","DB_PASS":"secret"}}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets, err := client.FetchSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error fetching secrets: %v", err)
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", secrets["API_KEY"])
	}
}

func TestWriteEnvFile_CreatesFile(t *testing.T) {
	tmpFile := t.TempDir() + "/.env"
	secrets := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	err := WriteEnvFile(tmpFile, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Error("expected non-empty env file")
	}
}
