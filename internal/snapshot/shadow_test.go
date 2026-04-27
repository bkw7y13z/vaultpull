package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeShadowSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snapshot.json")
	snap := Snapshot{Entries: entries}
	data, _ := json.MarshalIndent(snap, "", "  ")
	_ = os.WriteFile(path, data, 0o600)
	return path
}

func TestShadow_EmptyPath(t *testing.T) {
	err := Shadow("", "abc123", nil)
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestShadow_EmptyChecksum(t *testing.T) {
	path := writeShadowSnapshot(t, nil)
	err := Shadow(path, "", nil)
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestShadow_ChecksumNotFound(t *testing.T) {
	path := writeShadowSnapshot(t, []Entry{})
	err := Shadow(path, "notexist", nil)
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestShadow_Success(t *testing.T) {
	entry := Entry{
		Checksum:  "abc123",
		Tag:       "v1",
		Secrets:   map[string]string{"FOO": "bar", "BAZ": "qux"},
		CreatedAt: time.Now().UTC(),
	}
	path := writeShadowSnapshot(t, []Entry{entry})

	err := Shadow(path, "abc123", map[string]string{"env": "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	shadows, err := GetShadows(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error getting shadows: %v", err)
	}
	if len(shadows) != 1 {
		t.Fatalf("expected 1 shadow, got %d", len(shadows))
	}
	if shadows[0].Checksum != "abc123" {
		t.Errorf("expected checksum abc123, got %s", shadows[0].Checksum)
	}
	if shadows[0].Meta["env"] != "staging" {
		t.Errorf("expected meta env=staging, got %v", shadows[0].Meta)
	}
	if len(shadows[0].Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(shadows[0].Keys))
	}
}

func TestShadow_Appends(t *testing.T) {
	entry := Entry{
		Checksum: "sum1",
		Secrets:  map[string]string{"K": "v"},
	}
	path := writeShadowSnapshot(t, []Entry{entry})

	_ = Shadow(path, "sum1", map[string]string{"round": "1"})
	_ = Shadow(path, "sum1", map[string]string{"round": "2"})

	shadows, err := GetShadows(path, "sum1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(shadows) != 2 {
		t.Fatalf("expected 2 shadows, got %d", len(shadows))
	}
}

func TestGetShadows_EmptyPath(t *testing.T) {
	_, err := GetShadows("", "abc")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestGetShadows_EmptyChecksum(t *testing.T) {
	path := writeShadowSnapshot(t, nil)
	_, err := GetShadows(path, "")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}
