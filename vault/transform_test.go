package vault

import (
	"testing"
)

func TestApplyTransforms_Upper(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret"}
	rules := []TransformRule{{KeyPattern: "DB_PASS", Transform: "upper"}}
	out, err := ApplyTransforms(secrets, rules)
	if err != nil {
		t.Fatal(err)
	}
	if out["DB_PASS"] != "SECRET" {
		t.Errorf("expected SECRET, got %s", out["DB_PASS"])
	}
}

func TestApplyTransforms_WildcardTrim(t *testing.T) {
	secrets := map[string]string{"KEY_A": "  hello  ", "KEY_B": " world "}
	rules := []TransformRule{{KeyPattern: "KEY_*", Transform: "trim"}}
	out, err := ApplyTransforms(secrets, rules)
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY_A"] != "hello" || out["KEY_B"] != "world" {
		t.Errorf("unexpected values: %v", out)
	}
}

func TestApplyTransforms_Base64RoundTrip(t *testing.T) {
	secrets := map[string]string{"TOKEN": "mysecret"}
	rules := []TransformRule{{KeyPattern: "TOKEN", Transform: "base64"}}
	out, err := ApplyTransforms(secrets, rules)
	if err != nil {
		t.Fatal(err)
	}
	rules2 := []TransformRule{{KeyPattern: "TOKEN", Transform: "base64d"}}
	out2, err := ApplyTransforms(out, rules2)
	if err != nil {
		t.Fatal(err)
	}
	if out2["TOKEN"] != "mysecret" {
		t.Errorf("round trip failed: %s", out2["TOKEN"])
	}
}

func TestApplyTransforms_UnknownTransform(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	rules := []TransformRule{{KeyPattern: "KEY", Transform: "nonexistent"}}
	_, err := ApplyTransforms(secrets, rules)
	if err == nil {
		t.Error("expected error for unknown transform")
	}
}

func TestApplyTransforms_NoMatchingKeys(t *testing.T) {
	secrets := map[string]string{"OTHER": "val"}
	rules := []TransformRule{{KeyPattern: "DB_*", Transform: "upper"}}
	out, err := ApplyTransforms(secrets, rules)
	if err != nil {
		t.Fatal(err)
	}
	if out["OTHER"] != "val" {
		t.Errorf("expected unchanged value")
	}
}
