package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeVerifySnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(entries); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestVerify_EmptyPath(t *testing.T) {
	_, err := Verify("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestVerify_EmptySnapshot(t *testing.T) {
	path := writeVerifySnapshot(t, []Entry{})
	results, err := Verify(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestVerify_AllValid(t *testing.T) {
	entries := []Entry{
		{Checksum: "abc123", CreatedAt: time.Now()},
		{Checksum: "def456", CreatedAt: time.Now()},
	}
	path := writeVerifySnapshot(t, entries)
	results, err := Verify(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if !r.Valid {
			t.Errorf("expected valid, got reason: %s", r.Reason)
		}
	}
}

func TestVerify_MissingChecksum(t *testing.T) {
	entries := []Entry{
		{Checksum: "", CreatedAt: time.Now()},
	}
	path := writeVerifySnapshot(t, entries)
	results, err := Verify(path)
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Valid {
		t.Error("expected invalid")
	}
	if results[0].Reason != "missing checksum" {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestVerify_DuplicateChecksum(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Checksum: "dup999", CreatedAt: now},
		{Checksum: "dup999", CreatedAt: now.Add(time.Second)},
	}
	path := writeVerifySnapshot(t, entries)
	results, err := Verify(path)
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Valid == false {
		t.Error("first entry should be valid")
	}
	if results[1].Valid {
		t.Error("second entry should be invalid")
	}
	if results[1].Reason != "duplicate checksum" {
		t.Errorf("unexpected reason: %s", results[1].Reason)
	}
}
