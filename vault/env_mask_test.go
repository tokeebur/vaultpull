package vault

import (
	"testing"
)

func TestMaskSecrets_NoMatch(t *testing.T) {
	secrets := map[string]string{"HOST": "localhost", "PORT": "8080"}
	result := MaskSecrets(secrets, DefaultMaskRules())
	if result["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", result["HOST"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected 8080, got %s", result["PORT"])
	}
}

func TestMaskSecrets_FullMask(t *testing.T) {
	secrets := map[string]string{"DB_PASSWORD": "supersecret"}
	result := MaskSecrets(secrets, DefaultMaskRules())
	if result["DB_PASSWORD"] != "****" {
		t.Errorf("expected ****, got %s", result["DB_PASSWORD"])
	}
}

func TestMaskSecrets_ShowLast(t *testing.T) {
	secrets := map[string]string{"API_KEY": "abcdef1234"}
	result := MaskSecrets(secrets, DefaultMaskRules())
	if result["API_KEY"] != "****1234" {
		t.Errorf("expected ****1234, got %s", result["API_KEY"])
	}
}

func TestMaskSecrets_ShortValueShowLast(t *testing.T) {
	secrets := map[string]string{"MY_TOKEN": "ab"}
	result := MaskSecrets(secrets, DefaultMaskRules())
	if result["MY_TOKEN"] != "****" {
		t.Errorf("expected **** for short value, got %s", result["MY_TOKEN"])
	}
}

func TestMaskSecrets_CustomRule(t *testing.T) {
	rules := []MaskRule{{Key: "CUSTOM", ShowLast: 2}}
	secrets := map[string]string{"MY_CUSTOM_VAL": "hello"}
	result := MaskSecrets(secrets, rules)
	if result["MY_CUSTOM_VAL"] != "****lo" {
		t.Errorf("expected ****lo, got %s", result["MY_CUSTOM_VAL"])
	}
}

func TestMaskSecrets_EmptyMap(t *testing.T) {
	result := MaskSecrets(map[string]string{}, DefaultMaskRules())
	if len(result) != 0 {
		t.Errorf("expected empty map")
	}
}
