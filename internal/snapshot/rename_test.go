package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeRenameSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	snap := &Snapshot{Entries: entries}
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRenameTag_EmptyPath(t *testing.T) {
	_, err := RenameTag("", "old", "new")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestRenameTag_EmptyOldTag(t *testing.T) {
	_, err := RenameTag("/tmp/x.json", "", "new")
	if err == nil || err.Error() != "old tag is required" {
		t.Fatalf("expected old tag error, got %v", err)
	}
}

func TestRenameTag_EmptyNewTag(t *testing.T) {
	_, err := RenameTag("/tmp/x.json", "old", "")
	if err == nil || err.Error() != "new tag is required" {
		t.Fatalf("expected new tag error, got %v", err)
	}
}

func TestRenameTag_NotFound(t *testing.T) {
	path := writeRenameSnapshot(t, []Entry{
		{Checksum: "abc", Tag: "staging", Timestamp: time.Now()},
	})
	_, err := RenameTag(path, "production", "prod-v2")
	if err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestRenameTag_Success(t *testing.T) {
	path := writeRenameSnapshot(t, []Entry{
		{Checksum: "abc123", Tag: "v1", Timestamp: time.Now()},
	})
	res, err := RenameTag(path, "v1", "v1-stable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OldTag != "v1" || res.NewTag != "v1-stable" {
		t.Errorf("unexpected result: %+v", res)
	}

	snap, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if snap.Entries[0].Tag != "v1-stable" {
		t.Errorf("tag not persisted, got %q", snap.Entries[0].Tag)
	}
}
