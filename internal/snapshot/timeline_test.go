package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeTimelineSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(Store{Entries: entries}); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestTimeline_EmptyPath(t *testing.T) {
	_, err := Timeline("", TimelineOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTimeline_AllEntries(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	path := writeTimelineSnapshot(t, []Entry{
		{Checksum: "aaa", Keys: []string{"A"}, CreatedAt: now.Add(-2 * time.Hour)},
		{Checksum: "bbb", Keys: []string{"B", "C"}, CreatedAt: now.Add(-1 * time.Hour)},
		{Checksum: "ccc", Keys: []string{"D"}, CreatedAt: now},
	})
	entries, err := Timeline(path, TimelineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Checksum != "aaa" {
		t.Errorf("expected first entry aaa, got %s", entries[0].Checksum)
	}
}

func TestTimeline_SinceFilter(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	path := writeTimelineSnapshot(t, []Entry{
		{Checksum: "old", Keys: []string{"X"}, CreatedAt: now.Add(-5 * time.Hour)},
		{Checksum: "new", Keys: []string{"Y"}, CreatedAt: now},
	})
	entries, err := Timeline(path, TimelineOptions{Since: now.Add(-1 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Checksum != "new" {
		t.Errorf("expected only 'new', got %+v", entries)
	}
}

func TestTimeline_TaggedOnly(t *testing.T) {
	now := time.Now().UTC()
	path := writeTimelineSnapshot(t, []Entry{
		{Checksum: "t1", Tag: "v1", Keys: []string{"A"}, CreatedAt: now.Add(-time.Hour)},
		{Checksum: "t2", Keys: []string{"B"}, CreatedAt: now},
	})
	entries, err := Timeline(path, TimelineOptions{Tagged: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Tag != "v1" {
		t.Errorf("expected tagged only, got %+v", entries)
	}
}

func TestTimeline_KeyCount(t *testing.T) {
	now := time.Now().UTC()
	path := writeTimelineSnapshot(t, []Entry{
		{Checksum: "x", Keys: []string{"A", "B", "C"}, CreatedAt: now},
	})
	entries, err := Timeline(path, TimelineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].KeyCount != 3 {
		t.Errorf("expected KeyCount 3, got %d", entries[0].KeyCount)
	}
}
