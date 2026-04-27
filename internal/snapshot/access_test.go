package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeAccessSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	data := make(map[string]Entry)
	for _, e := range entries {
		data[e.Checksum] = e
	}
	b, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile(p, b, 0644)
	return p
}

func TestRecordAccess_EmptyPath(t *testing.T) {
	err := RecordAccess("", "abc", "alice", "read", "")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestRecordAccess_EmptyChecksum(t *testing.T) {
	p := writeAccessSnapshot(t, nil)
	err := RecordAccess(p, "", "alice", "read", "")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestRecordAccess_EmptyAccessedBy(t *testing.T) {
	p := writeAccessSnapshot(t, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	err := RecordAccess(p, "abc123", "", "read", "")
	if err == nil || err.Error() != "accessed_by is required" {
		t.Fatalf("expected accessed_by error, got %v", err)
	}
}

func TestRecordAccess_EmptyAction(t *testing.T) {
	p := writeAccessSnapshot(t, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	err := RecordAccess(p, "abc123", "alice", "", "")
	if err == nil || err.Error() != "action is required" {
		t.Fatalf("expected action error, got %v", err)
	}
}

func TestRecordAccess_ChecksumNotFound(t *testing.T) {
	p := writeAccessSnapshot(t, nil)
	err := RecordAccess(p, "missing", "alice", "read", "")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestRecordAccess_Success(t *testing.T) {
	p := writeAccessSnapshot(t, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	if err := RecordAccess(p, "abc123", "alice", "read", "audit"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, err := GetAccessLog(p, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].AccessedBy != "alice" || entries[0].Action != "read" || entries[0].Reason != "audit" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestGetAccessLog_Appends(t *testing.T) {
	p := writeAccessSnapshot(t, []Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	_ = RecordAccess(p, "abc123", "alice", "read", "")
	_ = RecordAccess(p, "abc123", "bob", "write", "deploy")
	entries, _ := GetAccessLog(p, "abc123")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestGetAccessLog_FiltersByChecksum(t *testing.T) {
	p := writeAccessSnapshot(t, []Entry{
		{Checksum: "abc123", CreatedAt: time.Now()},
		{Checksum: "def456", CreatedAt: time.Now()},
	})
	_ = RecordAccess(p, "abc123", "alice", "read", "")
	_ = RecordAccess(p, "def456", "bob", "read", "")
	entries, _ := GetAccessLog(p, "abc123")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry for abc123, got %d", len(entries))
	}
	if entries[0].Checksum != "abc123" {
		t.Errorf("unexpected checksum: %s", entries[0].Checksum)
	}
}
