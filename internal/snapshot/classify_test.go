package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeClassifySnapshot(t *testing.T, dir string) string {
	t.Helper()
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"DB_PASS"}, CapturedAt: time.Now()},
		},
	}
	data, _ := json.MarshalIndent(snap, "", "  ")
	path := filepath.Join(dir, "snap.json")
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestClassify_EmptyPath(t *testing.T) {
	err := Classify("", "abc123", "secret", "alice", true)
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestClassify_EmptyChecksum(t *testing.T) {
	dir := t.TempDir()
	path := writeClassifySnapshot(t, dir)
	err := Classify(path, "", "secret", "alice", true)
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestClassify_EmptyCategory(t *testing.T) {
	dir := t.TempDir()
	path := writeClassifySnapshot(t, dir)
	err := Classify(path, "abc123", "", "alice", false)
	if err == nil || err.Error() != "category is required" {
		t.Fatalf("expected category error, got %v", err)
	}
}

func TestClassify_EmptySetBy(t *testing.T) {
	dir := t.TempDir()
	path := writeClassifySnapshot(t, dir)
	err := Classify(path, "abc123", "secret", "", true)
	if err == nil || err.Error() != "set_by is required" {
		t.Fatalf("expected set_by error, got %v", err)
	}
}

func TestClassify_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	path := writeClassifySnapshot(t, dir)
	err := Classify(path, "notfound", "secret", "alice", false)
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestClassify_Success(t *testing.T) {
	dir := t.TempDir()
	path := writeClassifySnapshot(t, dir)
	err := Classify(path, "abc123", "confidential", "bob", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, ok, err := GetClassification(path, "abc123")
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	if !ok {
		t.Fatal("expected classification to be found")
	}
	if c.Category != "confidential" {
		t.Errorf("expected category confidential, got %s", c.Category)
	}
	if !c.Sensitive {
		t.Error("expected sensitive=true")
	}
	if c.SetBy != "bob" {
		t.Errorf("expected set_by=bob, got %s", c.SetBy)
	}
}

func TestGetClassification_NotFound(t *testing.T) {
	dir := t.TempDir()
	path := writeClassifySnapshot(t, dir)
	_, ok, err := GetClassification(path, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected not found")
	}
}
