package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeNoteSnapshot(t *testing.T, checksum string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	snap := &Snapshot{
		Entries: []Entry{
			{Checksum: checksum, CreatedAt: time.Now(), Keys: []string{"FOO"}},
		},
	}
	data, _ := encodeSnapshot(snap)
	os.WriteFile(p, data, 0600)
	return p
}

func TestAddNote_EmptyPath(t *testing.T) {
	err := AddNote("", "abc", "hello")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddNote_EmptyChecksum(t *testing.T) {
	err := AddNote("/tmp/snap.json", "", "hello")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddNote_EmptyNote(t *testing.T) {
	err := AddNote("/tmp/snap.json", "abc", "")
	if err == nil || err.Error() != "note is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddNote_ChecksumNotFound(t *testing.T) {
	p := writeNoteSnapshot(t, "abc123")
	err := AddNote(p, "notexist", "hello")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAddNote_Success(t *testing.T) {
	p := writeNoteSnapshot(t, "abc123")
	if err := AddNote(p, "abc123", "this is a note"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	note, err := GetNote(p, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note != "this is a note" {
		t.Errorf("expected 'this is a note', got %q", note)
	}
}

func TestGetNote_NotFound(t *testing.T) {
	p := writeNoteSnapshot(t, "abc123")
	_, err := GetNote(p, "abc123")
	if err == nil {
		t.Fatal("expected error for missing note")
	}
}

func TestAddNote_Overwrite(t *testing.T) {
	p := writeNoteSnapshot(t, "abc123")
	AddNote(p, "abc123", "first")
	AddNote(p, "abc123", "second")
	note, _ := GetNote(p, "abc123")
	if note != "second" {
		t.Errorf("expected 'second', got %q", note)
	}
}
