package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeBookmarkSnapshot(t *testing.T, checksum string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: checksum, CreatedAt: time.Now(), Keys: []string{"FOO"}},
		},
	}
	data, _ := json.Marshal(snap)
	os.WriteFile(path, data, 0644)
	return path
}

func TestAddBookmark_EmptyPath(t *testing.T) {
	err := AddBookmark("", "mylabel", "abc123", "")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddBookmark_EmptyLabel(t *testing.T) {
	err := AddBookmark("/tmp/snap.json", "", "abc123", "")
	if err == nil || err.Error() != "label is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddBookmark_EmptyChecksum(t *testing.T) {
	err := AddBookmark("/tmp/snap.json", "label", "", "")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddBookmark_ChecksumNotFound(t *testing.T) {
	path := writeBookmarkSnapshot(t, "known")
	err := AddBookmark(path, "label", "unknown", "")
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddBookmark_Success(t *testing.T) {
	path := writeBookmarkSnapshot(t, "abc123")
	if err := AddBookmark(path, "release-v1", "abc123", "initial release"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := GetBookmark(path, "release-v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Checksum != "abc123" {
		t.Errorf("expected checksum abc123, got %s", b.Checksum)
	}
	if b.Note != "initial release" {
		t.Errorf("expected note 'initial release', got %s", b.Note)
	}
}

func TestAddBookmark_DuplicateLabel(t *testing.T) {
	path := writeBookmarkSnapshot(t, "abc123")
	AddBookmark(path, "dup", "abc123", "")
	err := AddBookmark(path, "dup", "abc123", "")
	if err == nil || err.Error() != "bookmark label already exists" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestListBookmarks_Empty(t *testing.T) {
	path := writeBookmarkSnapshot(t, "abc123")
	list, err := ListBookmarks(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 bookmarks, got %d", len(list))
	}
}

func TestListBookmarks_Multiple(t *testing.T) {
	path := writeBookmarkSnapshot(t, "abc123")
	AddBookmark(path, "a", "abc123", "")
	AddBookmark(path, "b", "abc123", "")
	list, err := ListBookmarks(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 bookmarks, got %d", len(list))
	}
}
