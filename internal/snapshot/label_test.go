package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeLabelSnapshot(t *testing.T, checksum string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := &snapshot.Snapshot{
		Entries: []snapshot.Entry{
			{Checksum: checksum, CreatedAt: time.Now(), Keys: []string{"KEY"}},
		},
	}
	if err := snapshot.Save(path, snap); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLabel_EmptyPath(t *testing.T) {
	err := snapshot.Label("", "abc", "my-label")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestLabel_EmptyChecksum(t *testing.T) {
	err := snapshot.Label("/tmp/snap.json", "", "my-label")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestLabel_EmptyLabel(t *testing.T) {
	err := snapshot.Label("/tmp/snap.json", "abc", "")
	if err == nil || err.Error() != "label is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestLabel_ChecksumNotFound(t *testing.T) {
	path := writeLabelSnapshot(t, "aaa")
	err := snapshot.Label(path, "zzz", "my-label")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestLabel_Success(t *testing.T) {
	path := writeLabelSnapshot(t, "abc123")
	if err := snapshot.Label(path, "abc123", "production-baseline"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lbl, err := snapshot.GetLabel(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lbl != "production-baseline" {
		t.Errorf("expected 'production-baseline', got %q", lbl)
	}
}

func TestGetLabel_NotFound(t *testing.T) {
	path := writeLabelSnapshot(t, "abc123")
	_, err := snapshot.GetLabel(path, "abc123")
	if err == nil {
		t.Fatal("expected error for missing label")
	}
}

func TestLabel_StoreFileCreated(t *testing.T) {
	path := writeLabelSnapshot(t, "def456")
	if err := snapshot.Label(path, "def456", "staging"); err != nil {
		t.Fatal(err)
	}
	store := filepath.Join(filepath.Dir(path), "labels.json")
	if _, err := os.Stat(store); err != nil {
		t.Errorf("expected labels.json to exist: %v", err)
	}
}
