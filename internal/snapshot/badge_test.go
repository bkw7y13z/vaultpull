package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeBadgeSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	s := Snapshot{Entries: entries}
	if err := json.NewEncoder(f).Encode(s); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestSetBadge_EmptyPath(t *testing.T) {
	err := SetBadge("", "abc", "ci", "ok", "passed")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetBadge_EmptyChecksum(t *testing.T) {
	err := SetBadge("/tmp/snap.json", "", "ci", "ok", "passed")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetBadge_EmptyLabel(t *testing.T) {
	err := SetBadge("/tmp/snap.json", "abc", "", "ok", "passed")
	if err == nil || err.Error() != "label is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetBadge_InvalidStatus(t *testing.T) {
	path := writeBadgeSnapshot(t, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	err := SetBadge(path, "abc123", "ci", "unknown", "")
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestSetBadge_ChecksumNotFound(t *testing.T) {
	path := writeBadgeSnapshot(t, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	err := SetBadge(path, "notfound", "ci", "ok", "passed")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestSetBadge_Success(t *testing.T) {
	path := writeBadgeSnapshot(t, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	if err := SetBadge(path, "abc123", "ci", "ok", "all tests passed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, ok, err := GetBadge(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected badge to exist")
	}
	if b.Status != "ok" || b.Label != "ci" || b.Message != "all tests passed" {
		t.Fatalf("unexpected badge: %+v", b)
	}
}

func TestGetBadge_NotFound(t *testing.T) {
	path := writeBadgeSnapshot(t, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	_, ok, err := GetBadge(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected badge to not exist")
	}
}
