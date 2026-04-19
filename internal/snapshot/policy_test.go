package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writePolicySnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	s := Snapshot{Entries: entries}
	data, _ := json.Marshal(s)
	os.WriteFile(p, data, 0600)
	return p
}

func TestSetPolicy_EmptyPath(t *testing.T) {
	err := SetPolicy("", "abc", Policy{CreatedBy: "ci"})
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetPolicy_EmptyChecksum(t *testing.T) {
	p := writePolicySnapshot(t, nil)
	err := SetPolicy(p, "", Policy{CreatedBy: "ci"})
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetPolicy_EmptyCreatedBy(t *testing.T) {
	p := writePolicySnapshot(t, []Entry{{Checksum: "abc123", Timestamp: time.Now()}})
	err := SetPolicy(p, "abc123", Policy{})
	if err == nil || err.Error() != "created_by is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetPolicy_ChecksumNotFound(t *testing.T) {
	p := writePolicySnapshot(t, []Entry{{Checksum: "abc123", Timestamp: time.Now()}})
	err := SetPolicy(p, "notexist", Policy{CreatedBy: "ci"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSetPolicy_Success(t *testing.T) {
	p := writePolicySnapshot(t, []Entry{{Checksum: "abc123", Timestamp: time.Now()}})
	pol := Policy{MaxAge: 30, MinKeys: 5, RequireTag: true, CreatedBy: "ci"}
	if err := SetPolicy(p, "abc123", pol); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, found, err := GetPolicy(p, "abc123")
	if err != nil || !found {
		t.Fatalf("expected policy found, err=%v", err)
	}
	if got.MaxAge != 30 || got.MinKeys != 5 || !got.RequireTag {
		t.Errorf("unexpected policy: %+v", got)
	}
}

func TestSetPolicy_UpdateExisting(t *testing.T) {
	p := writePolicySnapshot(t, []Entry{{Checksum: "abc123", Timestamp: time.Now()}})
	SetPolicy(p, "abc123", Policy{MaxAge: 10, CreatedBy: "ci"})
	SetPolicy(p, "abc123", Policy{MaxAge: 60, CreatedBy: "admin"})
	got, _, _ := GetPolicy(p, "abc123")
	if got.MaxAge != 60 || got.CreatedBy != "admin" {
		t.Errorf("expected updated policy, got %+v", got)
	}
}

func TestGetPolicy_NotFound(t *testing.T) {
	p := writePolicySnapshot(t, []Entry{{Checksum: "abc123", Timestamp: time.Now()}})
	_, found, err := GetPolicy(p, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Error("expected not found")
	}
}
