package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeAnnotateSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "snap.json")
	snap := &Snapshot{Entries: entries}
	if err := Save(path, snap); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func TestAnnotate_EmptyPath(t *testing.T) {
	if err := Annotate("", "abc", "note"); err == nil {
		t.Fatal("expected error")
	}
}

func TestAnnotate_EmptyChecksum(t *testing.T) {
	if err := Annotate("/tmp/x.json", "", "note"); err == nil {
		t.Fatal("expected error")
	}
}

func TestAnnotate_EmptyNote(t *testing.T) {
	if err := Annotate("/tmp/x.json", "abc", ""); err == nil {
		t.Fatal("expected error")
	}
}

func TestAnnotate_ChecksumNotFound(t *testing.T) {
	path := writeAnnotateSnapshot(t, []Entry{
		{Checksum: "aaa", CapturedAt: time.Now()},
	})
	if err := Annotate(path, "bbb", "hello"); err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestAnnotate_Success(t *testing.T) {
	path := writeAnnotateSnapshot(t, []Entry{
		{Checksum: "abc123", CapturedAt: time.Now()},
	})
	if err := Annotate(path, "abc123", "deployed to prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	note, err := GetAnnotation(path, "abc123")
	if err != nil {
		t.Fatalf("get annotation: %v", err)
	}
	if note != "deployed to prod" {
		t.Errorf("expected %q, got %q", "deployed to prod", note)
	}
}

func TestGetAnnotation_NonExistentFile(t *testing.T) {
	_, err := GetAnnotation(filepath.Join(os.TempDir(), "no_such.json"), "abc")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAnnotation_ChecksumNotFound(t *testing.T) {
	path := writeAnnotateSnapshot(t, []Entry{
		{Checksum: "aaa", CapturedAt: time.Now()},
	})
	_, err := GetAnnotation(path, "zzz")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}
