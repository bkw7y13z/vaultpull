package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeLifecycleSnapshot(t *testing.T, dir string) string {
	t.Helper()
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"KEY"}, CreatedAt: time.Now()},
		},
	}
	data, _ := json.Marshal(snap)
	p := filepath.Join(dir, "snap.json")
	os.WriteFile(p, data, 0644)
	return p
}

func TestSetLifecycle_EmptyPath(t *testing.T) {
	err := SetLifecycle("", "abc", LifecycleActive, "user", "reason")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetLifecycle_EmptyChecksum(t *testing.T) {
	err := SetLifecycle("/tmp/x.json", "", LifecycleActive, "user", "reason")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetLifecycle_EmptyChangedBy(t *testing.T) {
	dir := t.TempDir()
	p := writeLifecycleSnapshot(t, dir)
	err := SetLifecycle(p, "abc123", LifecycleActive, "", "reason")
	if err == nil || err.Error() != "changed_by is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetLifecycle_EmptyReason(t *testing.T) {
	dir := t.TempDir()
	p := writeLifecycleSnapshot(t, dir)
	err := SetLifecycle(p, "abc123", LifecycleActive, "user", "")
	if err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetLifecycle_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	p := writeLifecycleSnapshot(t, dir)
	err := SetLifecycle(p, "notfound", LifecycleActive, "user", "reason")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSetLifecycle_Success(t *testing.T) {
	dir := t.TempDir()
	p := writeLifecycleSnapshot(t, dir)
	err := SetLifecycle(p, "abc123", LifecycleDeprecated, "admin", "old key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok, err := GetLifecycle(p, "abc123")
	if err != nil || !ok {
		t.Fatalf("expected entry, got err=%v ok=%v", err, ok)
	}
	if e.State != LifecycleDeprecated {
		t.Errorf("expected deprecated, got %s", e.State)
	}
	if e.ChangedBy != "admin" {
		t.Errorf("expected admin, got %s", e.ChangedBy)
	}
}

func TestGetLifecycle_NotFound(t *testing.T) {
	dir := t.TempDir()
	p := writeLifecycleSnapshot(t, dir)
	_, ok, err := GetLifecycle(p, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected not found")
	}
}

func TestSetLifecycle_UpdateExisting(t *testing.T) {
	dir := t.TempDir()
	p := writeLifecycleSnapshot(t, dir)
	SetLifecycle(p, "abc123", LifecycleActive, "user", "initial")
	SetLifecycle(p, "abc123", LifecycleRetired, "admin", "retiring")
	e, ok, err := GetLifecycle(p, "abc123")
	if err != nil || !ok {
		t.Fatalf("unexpected: err=%v ok=%v", err, ok)
	}
	if e.State != LifecycleRetired {
		t.Errorf("expected retired, got %s", e.State)
	}
}
