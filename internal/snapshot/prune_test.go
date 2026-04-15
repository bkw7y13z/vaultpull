package snapshot

import (
	"testing"
	"time"
)

func TestPrune_EmptyPath(t *testing.T) {
	_, err := Prune("", PruneOptions{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestPrune_NegativeKeepLast(t *testing.T) {
	_, err := Prune(tmpPath(t), PruneOptions{KeepLast: -1})
	if err == nil {
		t.Fatal("expected error for negative KeepLast")
	}
}

func TestPrune_NonExistentFile(t *testing.T) {
	// Load creates an empty snapshot when file is missing, so Prune should succeed.
	result, err := Prune(tmpPath(t), PruneOptions{KeepLast: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Removed != 0 || result.Retained != 0 {
		t.Errorf("expected 0/0, got %+v", result)
	}
}

func TestPrune_KeepLast(t *testing.T) {
	path := tmpPath(t)
	s, _ := Load(path)
	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		s.Entries = append(s.Entries, Entry{
			Timestamp: now.Add(time.Duration(i) * time.Minute),
			Checksum:  fmt.Sprintf("chk%d", i),
		})
	}
	_ = s.save(path)

	result, err := Prune(path, PruneOptions{KeepLast: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Removed != 2 || result.Retained != 3 {
		t.Errorf("expected removed=2 retained=3, got %+v", result)
	}
}

func TestPrune_MaxAge(t *testing.T) {
	path := tmpPath(t)
	s, _ := Load(path)
	now := time.Now().UTC()
	s.Entries = []Entry{
		{Timestamp: now.Add(-48 * time.Hour), Checksum: "old1"},
		{Timestamp: now.Add(-36 * time.Hour), Checksum: "old2"},
		{Timestamp: now.Add(-1 * time.Hour), Checksum: "recent"},
	}
	_ = s.save(path)

	result, err := Prune(path, PruneOptions{MaxAge: 24 * time.Hour})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Removed != 2 || result.Retained != 1 {
		t.Errorf("expected removed=2 retained=1, got %+v", result)
	}
}

func TestPrune_KeepLastOverridesAge(t *testing.T) {
	path := tmpPath(t)
	s, _ := Load(path)
	now := time.Now().UTC()
	s.Entries = []Entry{
		{Timestamp: now.Add(-72 * time.Hour), Checksum: "very-old"},
		{Timestamp: now.Add(-48 * time.Hour), Checksum: "old"},
	}
	_ = s.save(path)

	// Both entries are older than 1h, but KeepLast=1 should retain the newest.
	result, err := Prune(path, PruneOptions{MaxAge: time.Hour, KeepLast: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Removed != 1 || result.Retained != 1 {
		t.Errorf("expected removed=1 retained=1, got %+v", result)
	}
}
