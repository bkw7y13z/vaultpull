package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTagSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	data, _ := json.Marshal(entries)
	_ = os.WriteFile(p, data, 0600)
	return p
}

func TestTag_EmptyPath(t *testing.T) {
	err := Tag("", "abc123", "release-v1")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestTag_EmptyChecksum(t *testing.T) {
	err := Tag("/tmp/snap.json", "", "release-v1")
	if err == nil {
		t.Fatal("expected error for empty checksum")
	}
}

func TestTag_EmptyLabel(t *testing.T) {
	err := Tag("/tmp/snap.json", "abc123", "")
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestTag_ChecksumNotFound(t *testing.T) {
	p := writeTagSnapshot(t, []Entry{
		{Timestamp: time.Now(), Checksum: "aaa", Keys: []string{"KEY"}},
	})
	err := Tag(p, "nonexistent", "v1")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestTag_Success(t *testing.T) {
	p := writeTagSnapshot(t, []Entry{
		{Timestamp: time.Now(), Checksum: "abc123", Keys: []string{"DB_URL"}},
	})
	if err := Tag(p, "abc123", "release-v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := Load(p)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Tag != "release-v1" {
		t.Errorf("expected tag 'release-v1', got %q", entries[0].Tag)
	}
	if entries[0].TaggedAt.IsZero() {
		t.Error("expected TaggedAt to be set")
	}
}

func TestFindByTag_EmptyPath(t *testing.T) {
	_, err := FindByTag("", "v1")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestFindByTag_NotFound(t *testing.T) {
	p := writeTagSnapshot(t, []Entry{
		{Timestamp: time.Now(), Checksum: "abc", Tag: "other"},
	})
	_, err := FindByTag(p, "missing")
	if err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestFindByTag_Success(t *testing.T) {
	p := writeTagSnapshot(t, []Entry{
		{Timestamp: time.Now(), Checksum: "abc", Tag: "prod-deploy"},
	})
	e, err := FindByTag(p, "prod-deploy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Checksum != "abc" {
		t.Errorf("expected checksum 'abc', got %q", e.Checksum)
	}
}
