package vault

import (
	"testing"
	"time"
)

func TestSetAndIsExpired_NotExpired(t *testing.T) {
	m := make(TTLMap)
	SetTTL(m, "FOO", 10*time.Minute)
	if IsExpired(m, "FOO") {
		t.Fatal("expected FOO to not be expired")
	}
}

func TestIsExpired_Expired(t *testing.T) {
	m := make(TTLMap)
	m["BAR"] = TTLEntry{Key: "BAR", ExpiresAt: time.Now().Add(-1 * time.Second)}
	if !IsExpired(m, "BAR") {
		t.Fatal("expected BAR to be expired")
	}
}

func TestIsExpired_MissingKey(t *testing.T) {
	m := make(TTLMap)
	if IsExpired(m, "MISSING") {
		t.Fatal("missing key should not be expired")
	}
}

func TestExpiredKeys(t *testing.T) {
	m := make(TTLMap)
	SetTTL(m, "ALIVE", 10*time.Minute)
	m["DEAD"] = TTLEntry{Key: "DEAD", ExpiresAt: time.Now().Add(-1 * time.Second)}
	keys := ExpiredKeys(m)
	if len(keys) != 1 || keys[0] != "DEAD" {
		t.Fatalf("expected [DEAD], got %v", keys)
	}
}

func TestParseTTL_Duration(t *testing.T) {
	d, err := ParseTTL("5m")
	if err != nil || d != 5*time.Minute {
		t.Fatalf("expected 5m, got %v err=%v", d, err)
	}
}

func TestParseTTL_Seconds(t *testing.T) {
	d, err := ParseTTL("30")
	if err != nil || d != 30*time.Second {
		t.Fatalf("expected 30s, got %v err=%v", d, err)
	}
}

func TestParseTTL_Invalid(t *testing.T) {
	_, err := ParseTTL("abc")
	if err == nil {
		t.Fatal("expected error for invalid TTL")
	}
}

func TestFilterExpired_RemovesExpired(t *testing.T) {
	m := make(TTLMap)
	m["OLD"] = TTLEntry{Key: "OLD", ExpiresAt: time.Now().Add(-1 * time.Second)}
	SetTTL(m, "NEW", 10*time.Minute)
	secrets := map[string]string{"OLD": "val1", "NEW": "val2", "NOTAG": "val3"}
	result := FilterExpired(secrets, m)
	if _, ok := result["OLD"]; ok {
		t.Fatal("OLD should have been filtered")
	}
	if result["NEW"] != "val2" {
		t.Fatal("NEW should remain")
	}
	if result["NOTAG"] != "val3" {
		t.Fatal("NOTAG (no TTL) should remain")
	}
}
