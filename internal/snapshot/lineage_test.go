package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeLineageSnapshot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	snap := SnapshotFile{Entries: []Entry{
		{Checksum: "aaa111", Keys: []string{"FOO"}, CapturedAt: now()},
		{Checksum: "bbb222", Keys: []string{"BAR"}, CapturedAt: now()},
	}}
	data, _ := json.MarshalIndent(snap, "", "  ")
	path := filepath.Join(dir, "snapshot.json")
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestAddLineage_EmptyPath(t *testing.T) {
	err := AddLineage("", "aaa111", "bbb222", "derived")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestAddLineage_EmptyParent(t *testing.T) {
	path := writeLineageSnapshot(t)
	err := AddLineage(path, "", "bbb222", "derived")
	if err == nil || err.Error() != "parent checksum is required" {
		t.Fatalf("expected parent checksum error, got %v", err)
	}
}

func TestAddLineage_EmptyChild(t *testing.T) {
	path := writeLineageSnapshot(t)
	err := AddLineage(path, "aaa111", "", "derived")
	if err == nil || err.Error() != "child checksum is required" {
		t.Fatalf("expected child checksum error, got %v", err)
	}
}

func TestAddLineage_EmptyReason(t *testing.T) {
	path := writeLineageSnapshot(t)
	err := AddLineage(path, "aaa111", "bbb222", "")
	if err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestAddLineage_Success(t *testing.T) {
	path := writeLineageSnapshot(t)
	err := AddLineage(path, "aaa111", "bbb222", "promoted")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	edges, err := GetLineage(path, "bbb222")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
	if edges[0].Parent != "aaa111" {
		t.Errorf("expected parent aaa111, got %s", edges[0].Parent)
	}
	if edges[0].Reason != "promoted" {
		t.Errorf("expected reason 'promoted', got %s", edges[0].Reason)
	}
}

func TestGetLineage_EmptyPath(t *testing.T) {
	_, err := GetLineage("", "aaa111")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestGetLineage_EmptyChecksum(t *testing.T) {
	path := writeLineageSnapshot(t)
	_, err := GetLineage(path, "")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestGetLineage_NoEdges(t *testing.T) {
	path := writeLineageSnapshot(t)
	edges, err := GetLineage(path, "ccc333")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(edges))
	}
}
