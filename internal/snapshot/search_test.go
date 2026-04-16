package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeSearchSnapshot(t *testing.T, snap Snapshot) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestSearch_EmptyPath(t *testing.T) {
	_, err := Search("", SearchOptions{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSearch_NoFilters(t *testing.T) {
	now := time.Now()
	snap := Snapshot{Entries: []Entry{
		{Checksum: "aaa", Keys: []string{"FOO"}, CreatedAt: now},
		{Checksum: "bbb", Keys: []string{"BAR"}, CreatedAt: now},
	}}
	path := writeSearchSnapshot(t, snap)
	results, err := Search(path, SearchOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_ByKeyContains(t *testing.T) {
	now := time.Now()
	snap := Snapshot{Entries: []Entry{
		{Checksum: "aaa", Keys: []string{"DB_HOST"}, CreatedAt: now},
		{Checksum: "bbb", Keys: []string{"APP_SECRET"}, CreatedAt: now},
	}}
	path := writeSearchSnapshot(t, snap)
	results, err := Search(path, SearchOptions{KeyContains: "DB_"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Entry.Checksum != "aaa" {
		t.Fatalf("unexpected results: %+v", results)
	}
}

func TestSearch_ByTag(t *testing.T) {
	now := time.Now()
	snap := Snapshot{Entries: []Entry{
		{Checksum: "aaa", Tag: "v1", CreatedAt: now},
		{Checksum: "bbb", Tag: "v2", CreatedAt: now},
	}}
	path := writeSearchSnapshot(t, snap)
	results, err := Search(path, SearchOptions{TagEquals: "v1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Entry.Checksum != "aaa" {
		t.Fatalf("unexpected results: %+v", results)
	}
}

func TestSearch_ByTimeRange(t *testing.T) {
	base := time.Now()
	snap := Snapshot{Entries: []Entry{
		{Checksum: "old", CreatedAt: base.Add(-2 * time.Hour)},
		{Checksum: "new", CreatedAt: base.Add(2 * time.Hour)},
	}}
	path := writeSearchSnapshot(t, snap)
	results, err := Search(path, SearchOptions{After: base})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Entry.Checksum != "new" {
		t.Fatalf("unexpected results: %+v", results)
	}
}
