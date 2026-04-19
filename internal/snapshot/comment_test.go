package snapshot

import (
	"path/filepath"
	"testing"
	"time"
)

func writeCommentSnapshot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snapshot.json")
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"FOO"}, CreatedAt: time.Now()},
		},
	}
	if err := Save(path, snap); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func TestAddComment_EmptyPath(t *testing.T) {
	err := AddComment("", "abc", "hello", "alice")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddComment_EmptyChecksum(t *testing.T) {
	path := writeCommentSnapshot(t)
	err := AddComment(path, "", "hello", "alice")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddComment_EmptyText(t *testing.T) {
	path := writeCommentSnapshot(t)
	err := AddComment(path, "abc123", "", "alice")
	if err == nil || err.Error() != "comment text is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddComment_EmptyAuthor(t *testing.T) {
	path := writeCommentSnapshot(t)
	err := AddComment(path, "abc123", "hello", "")
	if err == nil || err.Error() != "author is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddComment_ChecksumNotFound(t *testing.T) {
	path := writeCommentSnapshot(t)
	err := AddComment(path, "notexist", "hello", "alice")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAddComment_Success(t *testing.T) {
	path := writeCommentSnapshot(t)
	if err := AddComment(path, "abc123", "looks good", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	comments, err := GetComments(path, "abc123")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
	if comments[0].Text != "looks good" || comments[0].Author != "alice" {
		t.Errorf("unexpected comment: %+v", comments[0])
	}
}

func TestGetComments_Empty(t *testing.T) {
	path := writeCommentSnapshot(t)
	comments, err := GetComments(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments")
	}
}
