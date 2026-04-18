package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func writeAuditSnapshot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snapshots.json")
	s := &SnapshotStore{}
	data, _ := marshalStore(s)
	os.WriteFile(p, data, 0600)
	return p
}

func TestRecordAudit_EmptyPath(t *testing.T) {
	err := RecordAudit("", "pull", "abc", "user", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRecordAudit_EmptyAction(t *testing.T) {
	p := writeAuditSnapshot(t)
	err := RecordAudit(p, "", "abc", "user", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRecordAudit_Success(t *testing.T) {
	p := writeAuditSnapshot(t)
	err := RecordAudit(p, "pull", "abc123", "alice", "synced 3 keys")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events, err := GetAuditLog(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Action != "pull" {
		t.Errorf("expected action 'pull', got %q", events[0].Action)
	}
	if events[0].Checksum != "abc123" {
		t.Errorf("expected checksum 'abc123', got %q", events[0].Checksum)
	}
}

func TestRecordAudit_Appends(t *testing.T) {
	p := writeAuditSnapshot(t)
	RecordAudit(p, "pull", "aaa", "alice", "")
	RecordAudit(p, "rotate", "bbb", "bob", "rotated")
	events, _ := GetAuditLog(p)
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestGetAuditLog_EmptyPath(t *testing.T) {
	_, err := GetAuditLog("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAuditLog_NonExistent(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "snapshots.json")
	events, err := GetAuditLog(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected empty log")
	}
}
