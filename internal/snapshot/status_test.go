package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeStatusSnapshot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := &Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"KEY"}, CreatedAt: time.Now()},
		},
	}
	data, _ := marshalSnapshot(snap)
	os.WriteFile(path, data, 0644)
	return path
}

func TestSetStatus_EmptyPath(t *testing.T) {
	err := SetStatus("", "abc123", "active", "user", "ok")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestSetStatus_EmptyChecksum(t *testing.T) {
	path := writeStatusSnapshot(t)
	err := SetStatus(path, "", "active", "user", "ok")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestSetStatus_InvalidState(t *testing.T) {
	path := writeStatusSnapshot(t)
	err := SetStatus(path, "abc123", "unknown", "user", "ok")
	if err == nil {
		t.Fatal("expected invalid state error")
	}
}

func TestSetStatus_ChecksumNotFound(t *testing.T) {
	path := writeStatusSnapshot(t)
	err := SetStatus(path, "notexist", "active", "user", "ok")
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestSetStatus_Success(t *testing.T) {
	path := writeStatusSnapshot(t)
	err := SetStatus(path, "abc123", "deprecated", "alice", "old version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry, err := GetStatus(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error on get: %v", err)
	}
	if entry.State != "deprecated" {
		t.Errorf("expected deprecated, got %s", entry.State)
	}
	if entry.SetBy != "alice" {
		t.Errorf("expected alice, got %s", entry.SetBy)
	}
}

func TestGetStatus_NotFound(t *testing.T) {
	path := writeStatusSnapshot(t)
	_, err := GetStatus(path, "abc123")
	if err == nil || err.Error() != "no status found for checksum" {
		t.Fatalf("expected not found error, got %v", err)
	}
}
