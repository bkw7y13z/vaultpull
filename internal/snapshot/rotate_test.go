package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRotate_EmptySnapshotPath(t *testing.T) {
	err := Rotate("", RotateOptions{ArchiveDir: t.TempDir()})
	if err == nil {
		t.Fatal("expected error for empty snapshot path")
	}
}

func TestRotate_EmptyArchiveDir(t *testing.T) {
	err := Rotate("/some/path", RotateOptions{})
	if err == nil {
		t.Fatal("expected error for empty archive dir")
	}
}

func TestRotate_NonExistentSnapshot(t *testing.T) {
	dir := t.TempDir()
	err := Rotate(filepath.Join(dir, "missing.json"), RotateOptions{ArchiveDir: dir})
	if err == nil {
		t.Fatal("expected error for missing snapshot file")
	}
}

func TestRotate_Success(t *testing.T) {
	src := t.TempDir()
	archive := t.TempDir()

	snapshotFile := filepath.Join(src, "snapshot.json")
	if err := os.WriteFile(snapshotFile, []byte(`{"entries":[]}`), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := Rotate(snapshotFile, RotateOptions{ArchiveDir: archive}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(archive)
	if err != nil {
		t.Fatalf("reading archive: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 archived file, got %d", len(entries))
	}
	name := entries[0].Name()
	if name[:8] != "snapshot" {
		t.Errorf("unexpected archive name prefix: %s", name)
	}
}

func TestRotate_PrunesOldArchives(t *testing.T) {
	src := t.TempDir()
	archive := t.TempDir()

	// Create a stale archive file with an old modification time.
	stale := filepath.Join(archive, "snapshot_old.json")
	if err := os.WriteFile(stale, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("setup stale: %v", err)
	}
	oldTime := time.Now().AddDate(0, 0, -10)
	if err := os.Chtimes(stale, oldTime, oldTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	snapshotFile := filepath.Join(src, "snapshot.json")
	if err := os.WriteFile(snapshotFile, []byte(`{"entries":[]}`), 0o644); err != nil {
		t.Fatalf("setup snapshot: %v", err)
	}

	if err := Rotate(snapshotFile, RotateOptions{ArchiveDir: archive, MaxAgeDays: 5}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := os.ReadDir(archive)
	for _, e := range entries {
		if e.Name() == "snapshot_old.json" {
			t.Error("stale archive file was not removed")
		}
	}
}
