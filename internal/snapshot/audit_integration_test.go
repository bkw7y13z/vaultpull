package snapshot

import (
	"path/filepath"
	"testing"
	"os"
)

func TestAudit_Integration_RecordAndRetrieve(t *testing.T) {
	dir := t.TempDir()
	snapshotPath := filepath.Join(dir, "snapshots.json")

	// Write a minimal snapshot store so auditPath resolves correctly
	s := &SnapshotStore{}
	data, err := marshalStore(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(snapshotPath, data, 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	actions := []struct{ action, checksum, actor, detail string }{
		{"pull", "aaa111", "alice", "initial sync"},
		{"rotate", "bbb222", "bob", "key rotation"},
		{"prune", "ccc333", "ci", "scheduled prune"},
	}

	for _, a := range actions {
		if err := RecordAudit(snapshotPath, a.action, a.checksum, a.actor, a.detail); err != nil {
			t.Fatalf("RecordAudit(%q): %v", a.action, err)
		}
	}

	events, err := GetAuditLog(snapshotPath)
	if err != nil {
		t.Fatalf("GetAuditLog: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	for i, a := range actions {
		if events[i].Action != a.action {
			t.Errorf("event[%d] action: got %q, want %q", i, events[i].Action, a.action)
		}
		if events[i].Actor != a.actor {
			t.Errorf("event[%d] actor: got %q, want %q", i, events[i].Actor, a.actor)
		}
		if events[i].Timestamp.IsZero() {
			t.Errorf("event[%d] timestamp is zero", i)
		}
	}
}
