package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFlagSnapshot(t *testing.T, dir string) string {
	t.Helper()
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"DB_PASS"}, CreatedAt: time.Now().UTC()},
		},
	}
	data, _ := json.Marshal(snap)
	p := filepath.Join(dir, "snap.json")
	os.WriteFile(p, data, 0644)
	return p
}

func TestFlag_EmptyPath(t *testing.T) {
	err := Flag("", "abc", "KEY", "val", "user", "reason")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFlag_EmptyChecksum(t *testing.T) {
	dir := t.TempDir()
	p := writeFlagSnapshot(t, dir)
	err := Flag(p, "", "KEY", "val", "user", "reason")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFlag_EmptyKey(t *testing.T) {
	dir := t.TempDir()
	p := writeFlagSnapshot(t, dir)
	err := Flag(p, "abc123", "", "val", "user", "reason")
	if err == nil || err.Error() != "key is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFlag_EmptyFlaggedBy(t *testing.T) {
	dir := t.TempDir()
	p := writeFlagSnapshot(t, dir)
	err := Flag(p, "abc123", "KEY", "val", "", "reason")
	if err == nil || err.Error() != "flagged_by is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFlag_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	p := writeFlagSnapshot(t, dir)
	err := Flag(p, "notexist", "KEY", "val", "user", "reason")
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestFlag_Success(t *testing.T) {
	dir := t.TempDir()
	p := writeFlagSnapshot(t, dir)
	err := Flag(p, "abc123", "DB_PASS", "secret", "alice", "suspicious value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	flags, err := GetFlags(p, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(flags) != 1 {
		t.Fatalf("expected 1 flag, got %d", len(flags))
	}
	if flags[0].Key != "DB_PASS" || flags[0].FlaggedBy != "alice" {
		t.Errorf("unexpected flag: %+v", flags[0])
	}
}

func TestGetFlags_AllWhenChecksumEmpty(t *testing.T) {
	dir := t.TempDir()
	p := writeFlagSnapshot(t, dir)
	_ = Flag(p, "abc123", "DB_PASS", "v1", "alice", "r1")
	flags, err := GetFlags(p, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(flags) != 1 {
		t.Errorf("expected 1, got %d", len(flags))
	}
}
