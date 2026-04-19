package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeLinkSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	s := Snapshot{Entries: entries}
	data, _ := json.Marshal(s)
	_ = os.WriteFile(p, data, 0644)
	return p
}

func TestAddLink_EmptyPath(t *testing.T) {
	err := AddLink("", "a", "b", "reason", "user")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddLink_EmptyFromChecksum(t *testing.T) {
	p := writeLinkSnapshot(t, nil)
	err := AddLink(p, "", "b", "reason", "user")
	if err == nil || err.Error() != "from checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddLink_EmptyToChecksum(t *testing.T) {
	p := writeLinkSnapshot(t, nil)
	err := AddLink(p, "a", "", "reason", "user")
	if err == nil || err.Error() != "to checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddLink_EmptyReason(t *testing.T) {
	p := writeLinkSnapshot(t, nil)
	err := AddLink(p, "a", "b", "", "user")
	if err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddLink_ChecksumNotFound(t *testing.T) {
	p := writeLinkSnapshot(t, []Entry{{Checksum: "aaa"}})
	err := AddLink(p, "missing", "aaa", "reason", "user")
	if err == nil {
		t.Fatal("expected error for missing from checksum")
	}
}

func TestAddLink_Success(t *testing.T) {
	entries := []Entry{
		{Checksum: "abc", Timestamp: time.Now()},
		{Checksum: "def", Timestamp: time.Now()},
	}
	p := writeLinkSnapshot(t, entries)
	if err := AddLink(p, "abc", "def", "supersedes", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	links, err := GetLinks(p, "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}
	if links[0].ToChecksum != "def" {
		t.Errorf("expected to=def, got %s", links[0].ToChecksum)
	}
	if links[0].Reason != "supersedes" {
		t.Errorf("expected reason=supersedes, got %s", links[0].Reason)
	}
}

func TestGetLinks_EmptyPath(t *testing.T) {
	_, err := GetLinks("", "abc")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetLinks_NoLinks(t *testing.T) {
	entries := []Entry{{Checksum: "abc", Timestamp: time.Now()}}
	p := writeLinkSnapshot(t, entries)
	links, err := GetLinks(p, "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected 0 links, got %d", len(links))
	}
}
