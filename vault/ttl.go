package vault

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TTLEntry holds a key and its expiry time.
type TTLEntry struct {
	Key       string
	ExpiresAt time.Time
}

// TTLMap maps secret keys to their expiry times.
type TTLMap map[string]TTLEntry

// SetTTL assigns a TTL duration to a key.
func SetTTL(m TTLMap, key string, ttl time.Duration) {
	m[key] = TTLEntry{
		Key:       key,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// IsExpired returns true if the key's TTL has elapsed.
func IsExpired(m TTLMap, key string) bool {
	entry, ok := m[key]
	if !ok {
		return false
	}
	return time.Now().After(entry.ExpiresAt)
}

// ExpiredKeys returns all keys that have passed their TTL.
func ExpiredKeys(m TTLMap) []string {
	var out []string
	for k, e := range m {
		if time.Now().After(e.ExpiresAt) {
			out = append(out, k)
		}
	}
	return out
}

// ParseTTL parses a human-readable TTL string like "30s", "5m", "2h".
func ParseTTL(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty TTL string")
	}
	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}
	// fallback: treat plain number as seconds
	if sec, err2 := strconv.Atoi(s); err2 == nil {
		return time.Duration(sec) * time.Second, nil
	}
	return 0, fmt.Errorf("invalid TTL %q: %w", s, err)
}

// FilterExpired removes expired keys from a secrets map using the TTLMap.
func FilterExpired(secrets map[string]string, m TTLMap) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if !IsExpired(m, k) {
			out[k] = v
		}
	}
	return out
}
