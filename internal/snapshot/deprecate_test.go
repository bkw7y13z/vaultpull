package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeDeprecateSnapshot(t *testing.T, dir string) string {
	t.Helper()
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"KEY"}, CreatedAt: time.Now()},
		},
	}
	data, _ := json.MarshalIndent(snap, "", "  ")
	p := filepath.Join(dir, "snap.json")
	_ = os.WriteFile(p, data, 0o644)
	return p
}

func TestDeprecate_EmptyPath(t *testing.T) {
	err := Deprecate("", "abc123", "outdated", "alice", "")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestDeprecate_EmptyChecksum(t *testing.T) {
	dir := t.TempDir()
	p := writeDeprecateSnapshot(t, dir)
	err := Deprecate(p, "", "outdated", "alice", "")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestDeprecate_EmptyReason(t *testing.T) {
	dir := t.TempDir()
	p := writeDeprecateSnapshot(t, dir)
	err := Deprecate(p, "abc123", "", "alice", "")
	if err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestDeprecate_EmptyDeprecatedBy(t *testing.T) {
	dir := t.TempDir()
	p := writeDeprecateSnapshot(t, dir)
	err := Deprecate(p, "abc123", "outdated", "", "")
	if err == nil || err.Error() != "deprecated_by is required" {
		t.Fatalf("expected deprecated_by error, got %v", err)
	}
}

func TestDeprecate_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	p := writeDeprecateSnapshot(t, dir)
	err := Deprecate(p, "nonexistent", "outdated", "alice", "")
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestDeprecate_Success(t *testing.T) {
	dir := t.TempDir()
	p := writeDeprecateSnapshot(t, dir)
	err := Deprecate(p, "abc123", "no longer used", "alice", "use xyz instead")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rec, err := GetDeprecation(p, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec == nil {
		t.Fatal("expected record, got nil")
	}
	if rec.Reason != "no longer used" {
		t.Errorf("expected reason 'no longer used', got %q", rec.Reason)
	}
	if rec.Suggest != "use xyz instead" {
		t.Errorf("expected suggest 'use xyz instead', got %q", rec.Suggest)
	}
	if rec.DeprecatedBy != "alice" {
		t.Errorf("expected deprecated_by 'alice', got %q", rec.DeprecatedBy)
	}
}

func TestGetDeprecation_NotFound(t *testing.T) {
	dir := t.TempDir()
	p := writeDeprecateSnapshot(t, dir)
	rec, err := GetDeprecation(p, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec != nil {
		t.Errorf("expected nil record for unknown checksum")
	}
}
