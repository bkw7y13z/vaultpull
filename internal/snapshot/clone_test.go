package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeCloneSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "snap.json")
	data, _ := json.Marshal(&Snapshot{Entries: entries})
	_ = os.WriteFile(p, data, 0600)
	return p
}

func TestClone_EmptyPath(t *testing.T) {
	_, err := Clone("", "abc", "")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestClone_EmptyRef(t *testing.T) {
	p := writeCloneSnapshot(t, nil)
	_, err := Clone(p, "", "")
	if err == nil || err.Error() != "ref is required" {
		t.Fatalf("expected ref error, got %v", err)
	}
}

func TestClone_RefNotFound(t *testing.T) {
	p := writeCloneSnapshot(t, []Entry{{Checksum: "aaa", Keys: []string{"K"}, CreatedAt: time.Now()}})
	_, err := Clone(p, "zzz", "")
	if err == nil {
		t.Fatal("expected error for missing ref")
	}
}

func TestClone_ByChecksum(t *testing.T) {
	p := writeCloneSnapshot(t, []Entry{
		{Checksum: "abc123", Keys: []string{"FOO", "BAR"}, CreatedAt: time.Now()},
	})
	res, err := Clone(p, "abc123", "my-clone")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SourceChecksum != "abc123" {
		t.Errorf("expected source abc123, got %s", res.SourceChecksum)
	}
	if res.Tag != "my-clone" {
		t.Errorf("expected tag my-clone, got %s", res.Tag)
	}
	snap, _ := Load(p)
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
}

func TestClone_ByTag(t *testing.T) {
	p := writeCloneSnapshot(t, []Entry{
		{Checksum: "def456", Tag: "release-1", Keys: []string{"SECRET"}, CreatedAt: time.Now()},
	})
	res, err := Clone(p, "release-1", "release-1-copy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SourceChecksum != "def456" {
		t.Errorf("unexpected source checksum: %s", res.SourceChecksum)
	}
}
