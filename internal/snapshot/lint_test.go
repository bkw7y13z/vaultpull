package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeLintSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(entries); err != nil {
		t.Fatalf("encode: %v", err)
	}
	return f.Name()
}

func TestLint_EmptyPath(t *testing.T) {
	_, err := Lint("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLint_CleanSnapshot(t *testing.T) {
	entries := []Entry{
		{Checksum: "abc123", Keys: []string{"KEY_A", "KEY_B"}, CreatedAt: time.Now()},
		{Checksum: "def456", Keys: []string{"KEY_C"}, CreatedAt: time.Now()},
	}
	path := writeLintSnapshot(t, entries)
	result, err := Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasIssues() {
		for _, i := range result.Issues {
			t.Logf("issue: %s", i)
		}
		t.Fatalf("expected no issues, got %d", len(result.Issues))
	}
}

func TestLint_DuplicateChecksum(t *testing.T) {
	entries := []Entry{
		{Checksum: "dup999", Keys: []string{"A"}, CreatedAt: time.Now()},
		{Checksum: "dup999", Keys: []string{"B"}, CreatedAt: time.Now()},
	}
	path := writeLintSnapshot(t, entries)
	result, err := Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasIssues() {
		t.Fatal("expected duplicate checksum issue")
	}
	if result.Issues[0].Field != "checksum" {
		t.Errorf("expected field=checksum, got %s", result.Issues[0].Field)
	}
}

func TestLint_MissingTimestamp(t *testing.T) {
	entries := []Entry{
		{Checksum: "ts001", Keys: []string{"X"}},
	}
	path := writeLintSnapshot(t, entries)
	result, err := Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, i := range result.Issues {
		if i.Field == "created_at" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected created_at issue")
	}
}

func TestLint_EmptyKeys(t *testing.T) {
	entries := []Entry{
		{Checksum: "ek001", Keys: []string{}, CreatedAt: time.Now()},
	}
	path := writeLintSnapshot(t, entries)
	result, err := Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasIssues() {
		t.Fatal("expected keys issue for empty keys slice")
	}
}

func TestLint_BlankKeyName(t *testing.T) {
	entries := []Entry{
		{Checksum: "bk001", Keys: []string{"VALID", "   "}, CreatedAt: time.Now()},
	}
	path := writeLintSnapshot(t, entries)
	result, err := Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasIssues() {
		t.Fatal("expected issue for blank key name")
	}
}
