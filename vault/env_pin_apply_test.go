package vault

import (
	"strings"
	"testing"
)

func TestApplyPinsWithReport_NoChange(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	out, reports := ApplyPinsWithReport(secrets, PinFile{})
	if len(reports) != 0 {
		t.Errorf("expected no reports")
	}
	if out["A"] != "1" {
		t.Errorf("value should be unchanged")
	}
}

func TestApplyPinsWithReport_Override(t *testing.T) {
	secrets := map[string]string{"DB": "old"}
	pf := PinFile{"DB": {Key: "DB", Value: "pinned"}}
	out, reports := ApplyPinsWithReport(secrets, pf)
	if out["DB"] != "pinned" {
		t.Errorf("expected pinned value")
	}
	if len(reports) != 1 {
		t.Fatalf("expected 1 report")
	}
	if reports[0].OldValue != "old" || reports[0].NewValue != "pinned" {
		t.Errorf("report values wrong")
	}
}

func TestApplyPinsWithReport_NewKey(t *testing.T) {
	secrets := map[string]string{}
	pf := PinFile{"NEW": {Key: "NEW", Value: "val"}}
	out, reports := ApplyPinsWithReport(secrets, pf)
	if out["NEW"] != "val" {
		t.Errorf("expected new key")
	}
	if reports[0].OldValue != "" {
		t.Errorf("old value should be empty for new key")
	}
}

func TestFormatPinReport_Empty(t *testing.T) {
	s := FormatPinReport(nil)
	if s != "No pins applied." {
		t.Errorf("unexpected: %q", s)
	}
}

func TestFormatPinReport_ShowsChanges(t *testing.T) {
	reports := []PinReport{
		{Key: "K", OldValue: "old", NewValue: "new", Pinned: true},
	}
	s := FormatPinReport(reports)
	if !strings.Contains(s, "[PIN]") {
		t.Errorf("expected [PIN] tag in output")
	}
	if !strings.Contains(s, "old") || !strings.Contains(s, "new") {
		t.Errorf("expected old and new values in output")
	}
}
