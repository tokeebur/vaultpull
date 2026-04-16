package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/vaultpull/vault"
)

func setupWatchServer(t *testing.T, secrets map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{"data": secrets},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
}

func TestWatchSecrets_DetectsChange(t *testing.T) {
	srv := setupWatchServer(t, map[string]interface{}{"KEY": "newvalue"})
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	dir := t.TempDir()
	outFile := filepath.Join(dir, ".env")

	// Write initial state
	if err := os.WriteFile(outFile, []byte("KEY=oldvalue\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	changed := make(chan vault.DiffResult, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opts := vault.WatchOptions{
		Client:     client,
		SecretPath: "secret/data/app",
		OutputFile: outFile,
		Interval:   200 * time.Millisecond,
		OnChange: func(diff vault.DiffResult) {
			changed <- diff
			cancel()
		},
	}

	go vault.WatchSecrets(ctx, opts) //nolint:errcheck

	select {
	case diff := <-changed:
		if diff.Changed["KEY"] != "newvalue" {
			t.Errorf("expected changed KEY=newvalue, got %v", diff.Changed)
		}
	case <-time.After(3 * time.Second):
		t.Error("timed out waiting for change detection")
	}
}

func TestWatchSecrets_NoChangeNoCallback(t *testing.T) {
	srv := setupWatchServer(t, map[string]interface{}{"KEY": "same"})
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	dir := t.TempDir()
	outFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(outFile, []byte("KEY=same\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	callbackFired := false
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	opts := vault.WatchOptions{
		Client:     client,
		SecretPath: "secret/data/app",
		OutputFile: outFile,
		Interval:   200 * time.Millisecond,
		OnChange:   func(_ vault.DiffResult) { callbackFired = true },
	}

	vault.WatchSecrets(ctx, opts) //nolint:errcheck

	if callbackFired {
		t.Error("expected no callback when secrets unchanged")
	}
}
