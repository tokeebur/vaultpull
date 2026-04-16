package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// vaultKVv2Response represents the KV v2 secret response from Vault.
type vaultKVv2Response struct {
	Data struct {
		Data map[string]string `json:"data"`
	} `json:"data"`
}

// FetchSecrets retrieves key-value pairs from a Vault KV v2 secret path.
func (c *Client) FetchSecrets(secretPath string) (map[string]string, error) {
	url := fmt.Sprintf("%s/v1/%s", c.addr, secretPath)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed: check your VAULT_TOKEN")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("secret path not found: %s", secretPath)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from Vault", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var result vaultKVv2Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing vault response: %w", err)
	}

	if result.Data.Data == nil {
		return nil, fmt.Errorf("no data found at path: %s", secretPath)
	}

	return result.Data.Data, nil
}
