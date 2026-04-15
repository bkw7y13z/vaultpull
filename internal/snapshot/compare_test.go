package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	s := &Snapshot{Entries: entries}
	if err := json.NewEncoder(f).Encode(s); err != nil {
		t.Fatalf("encode snapshot: %v", err)
	}
	return f.Name()
}

func TestCompare_EmptyPath(t *testing.T) {
	_, err := Compare("", "abc", "def")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestCompare_MissingChecksums(t *testing.T) {
	path := writeSnapshot(t, []Entry{})
	_, err := Compare(path, "", "def")
	if err == nil {
		t.Fatal("expected error for empty fromChecksum")
	}
	_, err = Compare(path, "abc", "")
	if err == nil {
		t.Fatal("expected error for empty toChecksum")
	}
}

func TestCompare_EntryNotFound(t *testing.T) {
	entries := []Entry{
		{Checksum: "aabbccdd", Timestamp: time.Now(), Keys: map[string]string{"A": "1"}},
	}
	path := writeSnapshot(t, entries)
	_, err := Compare(path, "aabb", "zzzz")
	if err == nil {
		t.Fatal("expected error for missing toChecksum")
	}
}

func TestCompare_Unchanged(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Checksum: "aaaa1111", Timestamp: now.Add(-time.Hour), Keys: map[string]string{"X": "1", "Y": "2"}},
		{Checksum: "bbbb2222", Timestamp: now, Keys: map[string]string{"X": "1", "Y": "2"}},
	}
	path := writeSnapshot(t, entries)
	res, err := Compare(path, "aaaa", "bbbb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Unchanged {
		t.Error("expected Unchanged=true")
	}
}

func TestCompare_WithChanges(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Checksum: "aaaa1111", Timestamp: now.Add(-time.Hour), Keys: map[string]string{"A": "1", "B": "old"}},
		{Checksum: "bbbb2222", Timestamp: now, Keys: map[string]string{"B": "new", "C": "3"}},
	}
	path := writeSnapshot(t, entries)
	res, err := Compare(path, "aaaa", "bbbb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Unchanged {
		t.Error("expected Unchanged=false")
	}
	if len(res.Diff.Added) != 1 || res.Diff.Added[0] != "C" {
		t.Errorf("expected Added=[C], got %v", res.Diff.Added)
	}
	if len(res.Diff.Removed) != 1 || res.Diff.Removed[0] != "A" {
		t.Errorf("expected Removed=[A], got %v", res.Diff.Removed)
	}
	if len(res.Diff.Changed) != 1 || res.Diff.Changed[0] != "B" {
		t.Errorf("expected Changed=[B], got %v", res.Diff.Changed)
	}
	summary := res.String()
	if summary == "" {
		t.Error("expected non-empty String() output")
	}
}
