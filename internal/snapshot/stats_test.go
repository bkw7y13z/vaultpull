package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeStatsSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "stats-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(entries); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestComputeStats_EmptyPath(t *testing.T) {
	_, err := ComputeStats("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestComputeStats_NonExistent(t *testing.T) {
	_, err := ComputeStats("/tmp/no-such-stats-file.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestComputeStats_EmptySnapshot(t *testing.T) {
	path := writeStatsSnapshot(t, []Entry{})
	s, err := ComputeStats(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.TotalEntries != 0 || s.UniqueKeys != 0 {
		t.Errorf("expected zeros, got %+v", s)
	}
}

func TestComputeStats_Counts(t *testing.T) {
	now := time.Now().UTC()
	entries := []Entry{
		{Checksum: "aaa", Keys: []string{"FOO", "BAR"}, Tag: "v1", Pinned: true, CreatedAt: now.Add(-time.Hour)},
		{Checksum: "bbb", Keys: []string{"FOO", "BAZ"}, Tag: "", Pinned: false, CreatedAt: now},
	}
	path := writeStatsSnapshot(t, entries)
	s, err := ComputeStats(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.TotalEntries != 2 {
		t.Errorf("expected 2 entries, got %d", s.TotalEntries)
	}
	if s.UniqueKeys != 3 {
		t.Errorf("expected 3 unique keys, got %d", s.UniqueKeys)
	}
	if s.TaggedEntries != 1 {
		t.Errorf("expected 1 tagged, got %d", s.TaggedEntries)
	}
	if s.PinnedEntries != 1 {
		t.Errorf("expected 1 pinned, got %d", s.PinnedEntries)
	}
	if len(s.TopKeys) == 0 || s.TopKeys[0] != "FOO" {
		t.Errorf("expected FOO as top key, got %v", s.TopKeys)
	}
	if !s.OldestAt.Before(s.NewestAt) {
		t.Errorf("expected oldest < newest")
	}
}
