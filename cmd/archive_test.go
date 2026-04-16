package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeArchiveSnap(t *testing.T, path string, entries []snapshot.Entry) {
	t.Helper()
	store := &snapshot.Store{Entries: entries}
	data, _ := json.MarshalIndent(store, "", "  ")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestRunArchive_MissingSnapshotFile(t *testing.T) {
	tmp := t.TempDir()
	cmd := archiveCmd
	cmd.Flags().Set("snapshot", filepath.Join(tmp, "missing.json"))
	cmd.Flags().Set("dir", filepath.Join(tmp, "arc"))
	cmd.Flags().Set("keep-last", "2")
	err := runArchive(cmd, nil)
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunArchive_Success(t *testing.T) {
	tmp := t.TempDir()
	snap := filepath.Join(tmp, "snap.json")
	arc := filepath.Join(tmp, "arc")
	now := time.Now()
	writeArchiveSnap(t, snap, []snapshot.Entry{
		{Checksum: "x1", CreatedAt: now.Add(-2 * time.Hour)},
		{Checksum: "x2", CreatedAt: now.Add(-1 * time.Hour)},
		{Checksum: "x3", CreatedAt: now},
	})

	cmd := archiveCmd
	cmd.Flags().Set("snapshot", snap)
	cmd.Flags().Set("dir", arc)
	cmd.Flags().Set("keep-last", "1")

	if err := runArchive(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(arc, "archive_index.json")); err != nil {
		t.Fatal("archive index not created")
	}
}

func TestArchiveCmd_RegisteredOnRoot(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "archive" {
			return
		}
	}
	t.Fatal("archive command not registered on root")
}

func TestArchiveCmd_DefaultFlags(t *testing.T) {
	f := archiveCmd.Flags()
	if v, _ := f.GetString("snapshot"); v != "snapshot.json" {
		t.Fatalf("unexpected default snapshot: %s", v)
	}
	if v, _ := f.GetInt("keep-last"); v != 10 {
		t.Fatalf("unexpected default keep-last: %d", v)
	}
}
