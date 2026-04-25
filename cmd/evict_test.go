package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"vaultpull/internal/snapshot"
)

func writeEvictSnap(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := snapshot.Snapshot{Entries: entries}
	data, _ := json.MarshalIndent(snap, "", "  ")
	_ = os.WriteFile(path, data, 0600)
	return path
}

func TestRunEvictAdd_MissingChecksum(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"evict", "add", "--by", "alice", "--reason", "stale"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "--checksum") {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestRunEvictAdd_MissingBy(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"evict", "add", "--checksum", "abc", "--reason", "stale"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "--by") {
		t.Fatalf("expected by error, got %v", err)
	}
}

func TestRunEvictAdd_Success(t *testing.T) {
	path := writeEvictSnap(t, []snapshot.Entry{
		{Checksum: "abc123", CreatedAt: time.Now()},
	})
	cmd := rootCmd
	cmd.SetArgs([]string{
		"evict", "add",
		"--snapshot", path,
		"--checksum", "abc123",
		"--by", "alice",
		"--reason", "expired",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	records, err := snapshot.GetEvictions(path)
	if err != nil {
		t.Fatalf("get evictions: %v", err)
	}
	if len(records) != 1 || records[0].Checksum != "abc123" {
		t.Fatalf("expected eviction record, got %+v", records)
	}
}

func TestRunEvictList_Empty(t *testing.T) {
	path := writeEvictSnap(t, nil)
	cmd := rootCmd
	cmd.SetArgs([]string{"evict", "list", "--snapshot", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEvictCmd_RegisteredOnRoot(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "evict" {
			return
		}
	}
	t.Fatal("evict command not registered on root")
}
