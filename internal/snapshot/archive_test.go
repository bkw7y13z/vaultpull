package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeArchiveSnapshot(t *testing.T, path string, entries []Entry) {
	t.Helper()
	store := &Store{Entries: entries}
	if err := Save(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
}

func TestArchive_EmptyPath(t *testing.T) {
	err := Archive("", "/tmp/arc", 2)
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestArchive_EmptyArchiveDir(t *testing.T) {
	err := Archive("/tmp/snap.json", "", 2)
	if err == nil || err.Error() != "archive dir is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestArchive_InvalidKeepLast(t *testing.T) {
	err := Archive("/tmp/snap.json", "/tmp/arc", 0)
	if err == nil {
		t.Fatal("expected error for keepLast=0")
	}
}

func TestArchive_NothingToArchive(t *testing.T) {
	tmp := t.TempDir()
	snap := filepath.Join(tmp, "snap.json")
	arc := filepath.Join(tmp, "arc")
	writeArchiveSnapshot(t, snap, []Entry{
		{Checksum: "aaa", CreatedAt: time.Now()},
	})
	if err := Archive(snap, arc, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// archive dir should not be created
	if _, err := os.Stat(arc); !os.IsNotExist(err) {
		t.Fatal("archive dir should not exist")
	}
}

func TestArchive_MovesOldEntries(t *testing.T) {
	tmp := t.TempDir()
	snap := filepath.Join(tmp, "snap.json")
	arc := filepath.Join(tmp, "arc")
	now := time.Now()
	entries := []Entry{
		{Checksum: "aaa", CreatedAt: now.Add(-3 * time.Hour)},
		{Checksum: "bbb", CreatedAt: now.Add(-2 * time.Hour)},
		{Checksum: "ccc", CreatedAt: now.Add(-1 * time.Hour)},
	}
	writeArchiveSnapshot(t, snap, entries)

	if err := Archive(snap, arc, 1); err != nil {
		t.Fatalf("archive: %v", err)
	}

	store, err := Load(snap)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(store.Entries) != 1 || store.Entries[0].Checksum != "ccc" {
		t.Fatalf("expected only ccc, got %+v", store.Entries)
	}

	idxPath := filepath.Join(arc, "archive_index.json")
	data, err := os.ReadFile(idxPath)
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	var idx ArchiveIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(idx.Entries) != 2 {
		t.Fatalf("expected 2 archive refs, got %d", len(idx.Entries))
	}
}
