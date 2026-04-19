package vault

import (
	"bytes"
	"strings"
	"testing"
)

func TestListGroups_Output(t *testing.T) {
	gm := GroupMap{
		"db":  {"DB_HOST", "DB_PORT"},
		"app": {"APP_KEY"},
	}
	var buf bytes.Buffer
	ListGroups(gm, &buf)
	out := buf.String()
	if !strings.Contains(out, "db") || !strings.Contains(out, "DB_HOST") {
		t.Fatalf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "app") {
		t.Fatalf("missing app group in output")
	}
}

func TestListGroups_Empty(t *testing.T) {
	var buf bytes.Buffer
	ListGroups(GroupMap{}, &buf)
	if !strings.Contains(buf.String(), "no groups") {
		t.Fatal("expected empty message")
	}
}

func TestGroupsForKey_Found(t *testing.T) {
	gm := GroupMap{
		"db":  {"DB_HOST", "DB_PORT"},
		"all": {"DB_HOST", "APP_KEY"},
	}
	groups := GroupsForKey(gm, "DB_HOST")
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %v", groups)
	}
}

func TestGroupsForKey_NotFound(t *testing.T) {
	gm := GroupMap{"db": {"DB_HOST"}}
	groups := GroupsForKey(gm, "MISSING")
	if len(groups) != 0 {
		t.Fatalf("expected no groups, got %v", groups)
	}
}
