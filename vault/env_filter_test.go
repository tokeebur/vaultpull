package vault

import (
	"testing"
)

var filterTestSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "vaultpull",
	"APP_VERSION": "1.0.0",
	"SECRET_KEY":  "abc123",
}

func TestFilterSecrets_Prefix(t *testing.T) {
	result, err := FilterSecrets(filterTestSecrets, FilterOptions{Mode: FilterModePrefix, Pattern: "DB_"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST")
	}
}

func TestFilterSecrets_Suffix(t *testing.T) {
	result, err := FilterSecrets(filterTestSecrets, FilterOptions{Mode: FilterModeSuffix, Pattern: "_KEY"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestFilterSecrets_Regex(t *testing.T) {
	result, err := FilterSecrets(filterTestSecrets, FilterOptions{Mode: FilterModeRegex, Pattern: "^APP_"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestFilterSecrets_InvalidRegex(t *testing.T) {
	_, err := FilterSecrets(filterTestSecrets, FilterOptions{Mode: FilterModeRegex, Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestFilterSecrets_Exact(t *testing.T) {
	result, err := FilterSecrets(filterTestSecrets, FilterOptions{Mode: FilterModeExact, Pattern: "APP_NAME"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestFilterSecrets_Invert(t *testing.T) {
	result, err := FilterSecrets(filterTestSecrets, FilterOptions{Mode: FilterModePrefix, Pattern: "DB_", Invert: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
}
