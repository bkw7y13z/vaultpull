package snapshot

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSummarize_EmptyPath(t *testing.T) {
	_, err := Summarize("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSummarize_NonExistentFile(t *testing.T) {
	s, err := Summarize(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalEntries != 0 {
		t.Errorf("expected 0 entries, got %d", s.TotalEntries)
	}
}

func TestSummarize_SingleEntry(t *testing.T) {
	path := tmpPath(t)

	entry := Entry{
		Checksum:  "abc123",
		Keys:      []string{"FOO", "BAR"},
		CreatedAt: time.Now().UTC(),
	}
	if err := save(path, []Entry{entry}); err != nil {
		t.Fatalf("save: %v", err)
	}

	s, err := Summarize(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalEntries != 1 {
		t.Errorf("expected 1 entry, got %d", s.TotalEntries)
	}
	if s.LatestChecksum != "abc123" {
		t.Errorf("expected checksum abc123, got %s", s.LatestChecksum)
	}
	if len(s.UniqueKeys) != 2 {
		t.Errorf("expected 2 unique keys, got %d", len(s.UniqueKeys))
	}
}

func TestSummarize_MultipleEntries_UniqueKeysMerged(t *testing.T) {
	path := tmpPath(t)

	now := time.Now().UTC()
	entries := []Entry{
		{Checksum: "c1", Keys: []string{"A", "B"}, CreatedAt: now.Add(-time.Hour)},
		{Checksum: "c2", Keys: []string{"B", "C"}, CreatedAt: now},
	}
	if err := save(path, entries); err != nil {
		t.Fatalf("save: %v", err)
	}

	s, err := Summarize(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalEntries != 2 {
		t.Errorf("expected 2 entries, got %d", s.TotalEntries)
	}
	if len(s.UniqueKeys) != 3 {
		t.Errorf("expected 3 unique keys (A,B,C), got %d: %v", len(s.UniqueKeys), s.UniqueKeys)
	}
	if s.LatestChecksum != "c2" {
		t.Errorf("expected latest checksum c2, got %s", s.LatestChecksum)
	}
}

func TestSummarize_MultipleEntries_OldestAndLatestAt(t *testing.T) {
	path := tmpPath(t)

	oldest := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	latest := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Checksum: "c1", Keys: []string{"A"}, CreatedAt: oldest},
		{Checksum: "c2", Keys: []string{"B"}, CreatedAt: latest},
	}
	if err := save(path, entries); err != nil {
		t.Fatalf("save: %v", err)
	}

	s, err := Summarize(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.OldestAt.Equal(oldest) {
		t.Errorf("expected oldest %v, got %v", oldest, s.OldestAt)
	}
	if !s.LatestAt.Equal(latest) {
		t.Errorf("expected latest %v, got %v", latest, s.LatestAt)
	}
}

func TestSummary_Print(t *testing.T) {
	s := &Summary{
		TotalEntries:   3,
		LatestChecksum: "deadbeef",
		LatestAt:       time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		OldestAt:       time.Date(2024, 5, 1, 8, 0, 0, 0, time.UTC),
		UniqueKeys:     []string{"KEY1", "KEY2"},
	}

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	for _, want := range []string{"3", "deadbeef", "2024-06-01", "2024-05-01", "2"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
