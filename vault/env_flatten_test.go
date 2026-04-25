package vault

import (
	"testing"
)

func TestFlattenSecrets_DotSeparator(t *testing.T) {
	secrets := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
	}
	out, err := FlattenSecrets(secrets, FlattenOptions{Separator: "_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %q", out["db_host"])
	}
	if out["db_port"] != "5432" {
		t.Errorf("expected db_port=5432, got %q", out["db_port"])
	}
}

func TestFlattenSecrets_Uppercase(t *testing.T) {
	secrets := map[string]string{
		"app.name": "vaultpull",
	}
	out, err := FlattenSecrets(secrets, FlattenOptions{Separator: "_", Uppercase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "vaultpull" {
		t.Errorf("expected APP_NAME=vaultpull, got %v", out)
	}
}

func TestFlattenSecrets_StripPrefix(t *testing.T) {
	secrets := map[string]string{
		"prod.db.host": "rds.example.com",
	}
	out, err := FlattenSecrets(secrets, FlattenOptions{Separator: "_", StripPrefix: "prod."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_host"] != "rds.example.com" {
		t.Errorf("expected db_host, got %v", out)
	}
}

func TestFlattenSecrets_Collision(t *testing.T) {
	secrets := map[string]string{
		"db.host": "localhost",
		"db_host": "remotehost",
	}
	_, err := FlattenSecrets(secrets, FlattenOptions{Separator: "_"})
	if err == nil {
		t.Error("expected collision error, got nil")
	}
}

func TestFlattenSecrets_NoSeparatorInKeys(t *testing.T) {
	secrets := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	out, err := FlattenSecrets(secrets, FlattenOptions{Separator: "_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" || out["PORT"] != "8080" {
		t.Errorf("expected unchanged keys, got %v", out)
	}
}

func TestFlattenSecrets_SlashSeparator(t *testing.T) {
	secrets := map[string]string{
		"infra/region": "us-east-1",
	}
	out, err := FlattenSecrets(secrets, FlattenOptions{Separator: "_", Uppercase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["INFRA_REGION"] != "us-east-1" {
		t.Errorf("expected INFRA_REGION, got %v", out)
	}
}

func TestFormatFlattenReport_ShowsRenames(t *testing.T) {
	original := map[string]string{"db.host": "localhost"}
	flattened := map[string]string{"db_host": "localhost"}
	report := FormatFlattenReport(original, flattened)
	if report == "" {
		t.Error("expected non-empty report")
	}
}
