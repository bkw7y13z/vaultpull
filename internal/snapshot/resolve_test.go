package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeResolveSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	data, _ := json.Marshal(entries)
	_ = os.WriteFile(p, data, 0o644)
	return p
}

func TestResolve_EmptyPath(t *testing.T) {
	_, err := Resolve("", "abc")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestResolve_EmptyRef(t *testing.T) {
	p := writeResolveSnapshot(t, nil)
	_, err := Resolve(p, "")
	if err == nil {
		t.Fatal("expected error for empty ref")
	}
}

func TestResolve_ByChecksum(t *testing.T) {
	entries := []Entry{
		{Checksum: "deadbeef", CreatedAt: time.Now()},
	}
	p := writeResolveSnapshot(t, entries)
	res, err := Resolve(p, "deadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Checksum != "deadbeef" {
		t.Errorf("expected deadbeef, got %s", res.Checksum)
	}
	if res.Method != "checksum" {
		t.Errorf("expected method=checksum, got %s", res.Method)
	}
}

func TestResolve_NotFound(t *testing.T) {
	entries := []Entry{
		{Checksum: "aabbccdd", CreatedAt: time.Now()},
	}
	p := writeResolveSnapshot(t, entries)
	_, err := Resolve(p, "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown ref")
	}
}

func TestResolve_NonExistentFile(t *testing.T) {
	_, err := Resolve("/tmp/does-not-exist-resolve.json", "abc")
	if err == nil {
		t.Fatal("expected error for missing snapshot file")
	}
}
