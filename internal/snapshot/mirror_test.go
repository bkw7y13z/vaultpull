package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeMirrorSnapshot(t *testing.T, path string) {
	t.Helper()
	entries := []snapshot.Entry{
		{Checksum: "abc123", Keys: []string{"FOO"}, CreatedAt: time.Now().UTC()},
	}
	if err := snapshot.Save(path, entries); err != nil {
		t.Fatalf("save: %v", err)
	}
}

func TestMirror_EmptyPath(t *testing.T) {
	err := snapshot.Mirror("", "abc", "/tmp/dest.json", "user")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestMirror_EmptyChecksum(t *testing.T) {
	err := snapshot.Mirror("/tmp/snap.json", "", "/tmp/dest.json", "user")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestMirror_EmptyDestPath(t *testing.T) {
	err := snapshot.Mirror("/tmp/snap.json", "abc", "", "user")
	if err == nil || err.Error() != "dest path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestMirror_EmptyMirroredBy(t *testing.T) {
	err := snapshot.Mirror("/tmp/snap.json", "abc", "/tmp/dest.json", "")
	if err == nil || err.Error() != "mirrored_by is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestMirror_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	snap := filepath.Join(dir, "snap.json")
	writeMirrorSnapshot(t, snap)

	err := snapshot.Mirror(snap, "notexist", filepath.Join(dir, "dest.json"), "user")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMirror_Success(t *testing.T) {
	dir := t.TempDir()
	snap := filepath.Join(dir, "snap.json")
	writeMirrorSnapshot(t, snap)

	dest := filepath.Join(dir, "mirror", "snap_copy.json")
	if err := snapshot.Mirror(snap, "abc123", dest, "ci-bot"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("dest file not created: %v", err)
	}

	idx := filepath.Join(dir, "mirror_index.json")
	if _, err := os.Stat(idx); err != nil {
		t.Fatalf("mirror index not created: %v", err)
	}
}
