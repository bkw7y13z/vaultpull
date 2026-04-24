package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeDetachSnapshot(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "snap.json")
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"KEY"}, CreatedAt: time.Now()},
		},
	}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestDetach_EmptyPath(t *testing.T) {
	err := Detach("", "abc123", "alice", "no longer needed")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestDetach_EmptyChecksum(t *testing.T) {
	dir := t.TempDir()
	path := writeDetachSnapshot(t, dir)
	err := Detach(path, "", "alice", "reason")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestDetach_EmptyDetachedBy(t *testing.T) {
	dir := t.TempDir()
	path := writeDetachSnapshot(t, dir)
	err := Detach(path, "abc123", "", "reason")
	if err == nil || err.Error() != "detached_by is required" {
		t.Fatalf("expected detached_by error, got %v", err)
	}
}

func TestDetach_EmptyReason(t *testing.T) {
	dir := t.TempDir()
	path := writeDetachSnapshot(t, dir)
	err := Detach(path, "abc123", "alice", "")
	if err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestDetach_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	path := writeDetachSnapshot(t, dir)
	err := Detach(path, "notexist", "alice", "reason")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestDetach_Success(t *testing.T) {
	dir := t.TempDir()
	path := writeDetachSnapshot(t, dir)

	err := Detach(path, "abc123", "alice", "migrating to new source")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rec, err := GetDetachment(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec == nil {
		t.Fatal("expected detach record, got nil")
	}
	if rec.DetachedBy != "alice" {
		t.Errorf("expected detached_by=alice, got %q", rec.DetachedBy)
	}
	if rec.Reason != "migrating to new source" {
		t.Errorf("unexpected reason: %q", rec.Reason)
	}
	if rec.DetachedAt.IsZero() {
		t.Error("expected non-zero DetachedAt")
	}
}

func TestGetDetachment_NotFound(t *testing.T) {
	dir := t.TempDir()
	path := writeDetachSnapshot(t, dir)

	rec, err := GetDetachment(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec != nil {
		t.Errorf("expected nil record, got %+v", rec)
	}
}
