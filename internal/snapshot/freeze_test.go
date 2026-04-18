package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFreezeSnapshot(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "snap.json")
	snap := &Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"FOO"}, CreatedAt: time.Now()},
		},
	}
	data, _ := json.Marshal(snap)
	os.WriteFile(path, data, 0644)
	return path
}

func TestFreeze_EmptyPath(t *testing.T) {
	err := Freeze("", "abc", "user", "reason")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFreeze_EmptyChecksum(t *testing.T) {
	err := Freeze("/tmp/snap.json", "", "user", "reason")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFreeze_EmptyFrozenBy(t *testing.T) {
	err := Freeze("/tmp/snap.json", "abc", "", "reason")
	if err == nil || err.Error() != "frozen_by is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFreeze_EmptyReason(t *testing.T) {
	err := Freeze("/tmp/snap.json", "abc", "user", "")
	if err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFreeze_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	path := writeFreezeSnapshot(t, dir)
	err := Freeze(path, "notexist", "user", "reason")
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestFreeze_Success(t *testing.T) {
	dir := t.TempDir()
	path := writeFreezeSnapshot(t, dir)
	if err := Freeze(path, "abc123", "alice", "compliance hold"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ok, err := IsFrozen(path, "abc123")
	if err != nil {
		t.Fatalf("IsFrozen error: %v", err)
	}
	if !ok {
		t.Fatal("expected entry to be frozen")
	}
}

func TestGetFreeze_RecordFields(t *testing.T) {
	dir := t.TempDir()
	path := writeFreezeSnapshot(t, dir)
	_ = Freeze(path, "abc123", "bob", "audit requirement")

	r, ok, err := GetFreeze(path, "abc123")
	if err != nil || !ok {
		t.Fatalf("expected record, got err=%v ok=%v", err, ok)
	}
	if r.FrozenBy != "bob" {
		t.Errorf("expected FrozenBy=bob, got %s", r.FrozenBy)
	}
	if r.Reason != "audit requirement" {
		t.Errorf("expected reason, got %s", r.Reason)
	}
}

func TestIsFrozen_NotFrozen(t *testing.T) {
	dir := t.TempDir()
	path := writeFreezeSnapshot(t, dir)
	ok, err := IsFrozen(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected not frozen")
	}
}
