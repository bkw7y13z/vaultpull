package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeRollbackSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "rollback-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(entries); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestRollback_EmptyPath(t *testing.T) {
	_, err := Rollback("", "abc")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestRollback_EmptyRef(t *testing.T) {
	_, err := Rollback("/tmp/x.json", "")
	if err == nil || err.Error() != "ref is required" {
		t.Fatalf("expected ref error, got %v", err)
	}
}

func TestRollback_EmptySnapshot(t *testing.T) {
	path := writeRollbackSnapshot(t, []Entry{})
	_, err := Rollback(path, "abc")
	if err == nil {
		t.Fatal("expected error for empty snapshot")
	}
}

func TestRollback_RefNotFound(t *testing.T) {
	entries := []Entry{
		{Checksum: "aaa111", CreatedAt: time.Now()},
	}
	path := writeRollbackSnapshot(t, entries)
	_, err := Rollback(path, "zzz999")
	if err == nil {
		t.Fatal("expected not-found error")
	}
}

func TestRollback_NoPreviousEntry(t *testing.T) {
	entries := []Entry{
		{Checksum: "aaa111", CreatedAt: time.Now()},
	}
	path := writeRollbackSnapshot(t, entries)
	_, err := Rollback(path, "aaa111")
	if err == nil {
		t.Fatal("expected error: no previous entry")
	}
}

func TestRollback_Success(t *testing.T) {
	entries := []Entry{
		{Checksum: "aaa111", CreatedAt: time.Now().Add(-2 * time.Hour)},
		{Checksum: "bbb222", CreatedAt: time.Now().Add(-1 * time.Hour)},
		{Checksum: "ccc333", CreatedAt: time.Now()},
	}
	path := writeRollbackSnapshot(t, entries)

	res, err := Rollback(path, "ccc333")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.RestoredChecksum != "bbb222" {
		t.Errorf("expected bbb222, got %s", res.RestoredChecksum)
	}
	if res.PreviousChecksum != "ccc333" {
		t.Errorf("expected ccc333, got %s", res.PreviousChecksum)
	}

	latest, err := Latest(path)
	if err != nil {
		t.Fatal(err)
	}
	if latest.Checksum != "bbb222" {
		t.Errorf("latest should be bbb222, got %s", latest.Checksum)
	}
}
