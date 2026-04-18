package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeLockSnapshot(t *testing.T, dir string, entries []Entry) string {
	t.Helper()
	path := filepath.Join(dir, "snap.json")
	snap := Snapshot{Entries: entries}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestLock_EmptyPath(t *testing.T) {
	err := Lock("", "abc", "alice", "freeze")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestLock_EmptyChecksum(t *testing.T) {
	err := Lock("/tmp/snap.json", "", "alice", "freeze")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestLock_EmptyLockedBy(t *testing.T) {
	err := Lock("/tmp/snap.json", "abc", "", "freeze")
	if err == nil || err.Error() != "locked_by is required" {
		t.Fatalf("expected locked_by error, got %v", err)
	}
}

func TestLock_EmptyReason(t *testing.T) {
	err := Lock("/tmp/snap.json", "abc", "alice", "")
	if err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestLock_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	path := writeLockSnapshot(t, dir, []Entry{{Checksum: "aaa", CreatedAt: time.Now()}})
	err := Lock(path, "zzz", "alice", "freeze")
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestLock_Success(t *testing.T) {
	dir := t.TempDir()
	path := writeLockSnapshot(t, dir, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	if err := Lock(path, "abc123", "alice", "production freeze"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	locked, entry, err := IsLocked(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !locked {
		t.Fatal("expected entry to be locked")
	}
	if entry.LockedBy != "alice" {
		t.Errorf("expected alice, got %s", entry.LockedBy)
	}
}

func TestLock_AlreadyLocked(t *testing.T) {
	dir := t.TempDir()
	path := writeLockSnapshot(t, dir, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	_ = Lock(path, "abc123", "alice", "freeze")
	err := Lock(path, "abc123", "bob", "another freeze")
	if err == nil || err.Error() != "entry is already locked" {
		t.Fatalf("expected already locked error, got %v", err)
	}
}

func TestUnlock_Success(t *testing.T) {
	dir := t.TempDir()
	path := writeLockSnapshot(t, dir, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	_ = Lock(path, "abc123", "alice", "freeze")
	if err := Unlock(path, "abc123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	locked, _, _ := IsLocked(path, "abc123")
	if locked {
		t.Fatal("expected entry to be unlocked")
	}
}

func TestUnlock_NotFound(t *testing.T) {
	dir := t.TempDir()
	path := writeLockSnapshot(t, dir, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	err := Unlock(path, "zzz")
	if err == nil || err.Error() != "lock not found" {
		t.Fatalf("expected lock not found error, got %v", err)
	}
}
