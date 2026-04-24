package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeImpactSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}

func TestImpact_EmptyPath(t *testing.T) {
	_, err := Impact("", "DB_PASSWORD")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestImpact_EmptyKey(t *testing.T) {
	path := writeImpactSnapshot(t, []Entry{})
	_, err := Impact(path, "")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestImpact_EmptySnapshot(t *testing.T) {
	path := writeImpactSnapshot(t, []Entry{})
	report, err := Impact(path, "DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.TotalRefs != 0 {
		t.Errorf("expected 0 refs, got %d", report.TotalRefs)
	}
}

func TestImpact_KeyFound(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	entries := []Entry{
		{Checksum: "abc123", Keys: []string{"DB_PASSWORD", "API_KEY"}, CreatedAt: now},
		{Checksum: "def456", Keys: []string{"DB_PASSWORD"}, CreatedAt: now.Add(-time.Hour)},
		{Checksum: "ghi789", Keys: []string{"API_KEY"}, CreatedAt: now.Add(-2 * time.Hour)},
	}
	path := writeImpactSnapshot(t, entries)
	report, err := Impact(path, "DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.TotalRefs != 2 {
		t.Errorf("expected 2 refs, got %d", report.TotalRefs)
	}
	if report.Key != "DB_PASSWORD" {
		t.Errorf("expected key DB_PASSWORD, got %s", report.Key)
	}
}

func TestImpact_KeyNotPresent(t *testing.T) {
	now := time.Now().UTC()
	entries := []Entry{
		{Checksum: "abc123", Keys: []string{"API_KEY"}, CreatedAt: now},
	}
	path := writeImpactSnapshot(t, entries)
	report, err := Impact(path, "DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.TotalRefs != 0 {
		t.Errorf("expected 0 refs, got %d", report.TotalRefs)
	}
	if len(report.Entries) != 0 {
		t.Errorf("expected no entries, got %d", len(report.Entries))
	}
}
