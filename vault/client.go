package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// Client is a minimal Vault HTTP client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Vault client.
func NewClient(baseURL, token string) (*Client, error) {
	if baseURL == "" || token == "" {
		return nil, fmt.Errorf("baseURL and token are required")
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// GetSecrets fetches key-value pairs from a KV v2 secret path.
func (c *Client) GetSecrets(path string) (map[string]string, error) {
	url := fmt.Sprintf("%s/v1/%s", c.baseURL, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault returned status %d for path %s", resp.StatusCode, path)
	}

	var {
		Data struct {
			Data map[string]interfacet} `json:"data"`
	}
	if err := json.NewDecode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode vault response: %w", err)
	}

	secrets := make(map[string]string, len(result.Data.Data))
	for k, v := range result.Data.Data {
		secrets[k] = fmt.Sprintf("%v", v)
	}
	return secrets, nil
}

// WriteEnvFile writes secrets to a .env file in KEY=VALUE format.
func WriteEnvFile(path string, secrets map[string]string) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%q\n", k, secrets[k])
	}

	return os.WriteFile(path, []byte(sb.String()), 0600)
}
