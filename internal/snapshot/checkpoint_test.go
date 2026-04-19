package snapshot

import (
	"testing"
	"time"
)

func writeCheckpointSnapshot(t *testing.T) string {
	t.Helper()
	p := tmpPath(t)
	entry := Entry{
		Checksum:  "abc123",
		Keys:      []string{"KEY1"},
		CreatedAt: time.Now().UTC(),
	}
	if err := Save(p, []Entry{entry}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return p
}

func TestSetCheckpoint_EmptyPath(t *testing.T) {
	err := SetCheckpoint("", "abc123", "release-1", "alice")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetCheckpoint_EmptyChecksum(t *testing.T) {
	p := writeCheckpointSnapshot(t)
	err := SetCheckpoint(p, "", "release-1", "alice")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetCheckpoint_EmptyLabel(t *testing.T) {
	p := writeCheckpointSnapshot(t)
	err := SetCheckpoint(p, "abc123", "", "alice")
	if err == nil || err.Error() != "label is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetCheckpoint_ChecksumNotFound(t *testing.T) {
	p := writeCheckpointSnapshot(t)
	err := SetCheckpoint(p, "notfound", "release-1", "alice")
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSetCheckpoint_Success(t *testing.T) {
	p := writeCheckpointSnapshot(t)
	if err := SetCheckpoint(p, "abc123", "release-1", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cps, err := GetCheckpoints(p, "abc123")
	if err != nil {
		t.Fatalf("GetCheckpoints: %v", err)
	}
	if len(cps) != 1 {
		t.Fatalf("expected 1 checkpoint, got %d", len(cps))
	}
	if cps[0].Label != "release-1" {
		t.Errorf("expected label release-1, got %s", cps[0].Label)
	}
	if cps[0].CreatedBy != "alice" {
		t.Errorf("expected createdBy alice, got %s", cps[0].CreatedBy)
	}
}

func TestSetCheckpoint_Appends(t *testing.T) {
	p := writeCheckpointSnapshot(t)
	_ = SetCheckpoint(p, "abc123", "v1", "alice")
	_ = SetCheckpoint(p, "abc123", "v2", "bob")
	cps, _ := GetCheckpoints(p, "abc123")
	if len(cps) != 2 {
		t.Fatalf("expected 2 checkpoints, got %d", len(cps))
	}
}

func TestGetCheckpoints_EmptyPath(t *testing.T) {
	_, err := GetCheckpoints("", "abc123")
	if err == nil {
		t.Fatal("expected error")
	}
}
