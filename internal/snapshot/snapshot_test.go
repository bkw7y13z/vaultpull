package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicholasgasior/vaultpull/internal/snapshot"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestLoad_NonExistent(t *testing.T) {
	s, err := snapshot.Load("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(s.Entries))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tmpPath(t)
	s := &snapshot.Snapshot{}
	s.Add(snapshot.Entry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Namespace: "prod",
		Keys:      []string{"DB_HOST", "DB_PASS"},
		Checksum:  "abc123",
	}, 10)

	if err := s.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Namespace != "prod" {
		t.Errorf("expected namespace 'prod', got %q", loaded.Entries[0].Namespace)
	}
}

func TestAdd_MaxEntries(t *testing.T) {
	s := &snapshot.Snapshot{}
	for i := 0; i < 5; i++ {
		s.Add(snapshot.Entry{Namespace: "ns"}, 3)
	}
	if len(s.Entries) != 3 {
		t.Errorf("expected 3 entries after max trim, got %d", len(s.Entries))
	}
}

func TestLatest_Empty(t *testing.T) {
	s := &snapshot.Snapshot{}
	if s.Latest() != nil {
		t.Error("expected nil for empty snapshot")
	}
}

func TestLatest_ReturnsLast(t *testing.T) {
	s := &snapshot.Snapshot{}
	s.Add(snapshot.Entry{Namespace: "first"}, 10)
	s.Add(snapshot.Entry{Namespace: "last"}, 10)
	if s.Latest().Namespace != "last" {
		t.Errorf("expected 'last', got %q", s.Latest().Namespace)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tmpPath(t)
	os.WriteFile(path, []byte("not-json{"), 0600)
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
