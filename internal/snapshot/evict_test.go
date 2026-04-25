package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeEvictSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := Snapshot{Entries: entries}
	data, _ := json.MarshalIndent(snap, "", "  ")
	_ = os.WriteFile(path, data, 0600)
	return path
}

func TestEvict_EmptyPath(t *testing.T) {
	err := Evict("", "abc", "user", "stale")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestEvict_EmptyChecksum(t *testing.T) {
	err := Evict("/tmp/snap.json", "", "user", "stale")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestEvict_EmptyEvictedBy(t *testing.T) {
	err := Evict("/tmp/snap.json", "abc", "", "stale")
	if err == nil || err.Error() != "evicted_by is required" {
		t.Fatalf("expected evicted_by error, got %v", err)
	}
}

func TestEvict_EmptyReason(t *testing.T) {
	err := Evict("/tmp/snap.json", "abc", "user", "")
	if err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestEvict_ChecksumNotFound(t *testing.T) {
	path := writeEvictSnapshot(t, []Entry{
		{Checksum: "aaa", CreatedAt: time.Now()},
	})
	err := Evict(path, "zzz", "user", "stale")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestEvict_Success(t *testing.T) {
	path := writeEvictSnapshot(t, []Entry{
		{Checksum: "aaa", CreatedAt: time.Now()},
		{Checksum: "bbb", CreatedAt: time.Now()},
	})
	if err := Evict(path, "aaa", "alice", "no longer needed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(snap.Entries) != 1 || snap.Entries[0].Checksum != "bbb" {
		t.Fatalf("expected only bbb remaining, got %+v", snap.Entries)
	}

	records, err := GetEvictions(path)
	if err != nil {
		t.Fatalf("get evictions: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 eviction record, got %d", len(records))
	}
	if records[0].Checksum != "aaa" || records[0].EvictedBy != "alice" {
		t.Fatalf("unexpected record: %+v", records[0])
	}
}

func TestGetEvictions_EmptyPath(t *testing.T) {
	_, err := GetEvictions("")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestGetEvictions_NoFile(t *testing.T) {
	path := writeEvictSnapshot(t, nil)
	records, err := GetEvictions(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty, got %d", len(records))
	}
}
