package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"vaultpull/internal/snapshot"
)

func writeBookmarkSnap(t *testing.T, checksum string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := snapshot.Snapshot{
		Entries: []snapshot.Entry{
			{Checksum: checksum, CreatedAt: time.Now(), Keys: []string{"KEY"}},
		},
	}
	data, _ := json.Marshal(snap)
	os.WriteFile(path, data, 0644)
	return path
}

func TestRunBookmarkAdd_MissingLabel(t *testing.T) {
	cmd := bookmarkAddCmd
	cmd.Flags().Set("label", "")
	cmd.Flags().Set("checksum", "abc")
	err := runBookmarkAdd(cmd, nil)
	if err == nil || err.Error() != "--label is required" {
		t.Fatalf("expected label error, got %v", err)
	}
}

func TestRunBookmarkAdd_MissingChecksum(t *testing.T) {
	cmd := bookmarkAddCmd
	cmd.Flags().Set("label", "v1")
	cmd.Flags().Set("checksum", "")
	err := runBookmarkAdd(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestRunBookmarkAdd_Success(t *testing.T) {
	path := writeBookmarkSnap(t, "deadbeef")
	cmd := bookmarkAddCmd
	cmd.Flags().Set("snapshot", path)
	cmd.Flags().Set("label", "prod")
	cmd.Flags().Set("checksum", "deadbeef")
	cmd.Flags().Set("note", "production release")
	if err := runBookmarkAdd(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunBookmarkList_Empty(t *testing.T) {
	path := writeBookmarkSnap(t, "deadbeef")
	cmd := bookmarkListCmd
	cmd.Flags().Set("snapshot", path)
	if err := runBookmarkList(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookmarkCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "bookmark" {
			found = true
			break
		}
	}
	if !found {
		t.Error("bookmark command not registered on root")
	}
}
