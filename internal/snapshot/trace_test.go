package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTraceSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snapshot.json")
	data, _ := json.Marshal(entries)
	os.WriteFile(p, data, 0644)
	return p
}

func TestAddTrace_EmptyPath(t *testing.T) {
	err := AddTrace("", "abc", "pull", "user", "")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddTrace_EmptyChecksum(t *testing.T) {
	err := AddTrace("/tmp/s.json", "", "pull", "user", "")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddTrace_EmptyOperation(t *testing.T) {
	err := AddTrace("/tmp/s.json", "abc", "", "user", "")
	if err == nil || err.Error() != "operation is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddTrace_EmptyActor(t *testing.T) {
	err := AddTrace("/tmp/s.json", "abc", "pull", "", "")
	if err == nil || err.Error() != "actor is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddTrace_ChecksumNotFound(t *testing.T) {
	p := writeTraceSnapshot(t, []Entry{{Checksum: "aaa", Keys: []string{"K"}}})
	err := AddTrace(p, "zzz", "pull", "user", "")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestAddTrace_Success(t *testing.T) {
	p := writeTraceSnapshot(t, []Entry{{Checksum: "abc123", Keys: []string{"KEY"}, Timestamp: time.Now()}})
	err := AddTrace(p, "abc123", "pull", "alice", "initial sync")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events, err := GetTraces(p, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Operation != "pull" || events[0].Actor != "alice" || events[0].Detail != "initial sync" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestGetTraces_EmptyPath(t *testing.T) {
	_, err := GetTraces("", "abc")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetTraces_EmptyChecksum(t *testing.T) {
	_, err := GetTraces("/tmp/s.json", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetTraces_Appends(t *testing.T) {
	p := writeTraceSnapshot(t, []Entry{{Checksum: "abc123", Keys: []string{"KEY"}, Timestamp: time.Now()}})
	AddTrace(p, "abc123", "pull", "alice", "")
	AddTrace(p, "abc123", "rotate", "bob", "scheduled")
	events, err := GetTraces(p, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}
