package cmd

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeGCSnap(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "gc-*.json")
	if err != nil {
		t.Fatal(err)
	}
	snap := snapshot.Snapshot{Entries: entries}
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRunGC_MissingSnapshot(t *testing.T) {
	rootCmd.SetArgs([]string{"gc", "--snapshot", "/nonexistent/snap.json", "--max-age", "1h"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRunGC_Success(t *testing.T) {
	old := snapshot.Entry{Checksum: "abc", CreatedAt: time.Now().Add(-48 * time.Hour)}
	path := writeGCSnap(t, []snapshot.Entry{old})

	rootCmd.SetArgs([]string{"gc", "--snapshot", path, "--max-age", "24h",
		"--keep-pinned=false", "--keep-tagged=false"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGCCmd_RegisteredOnRoot(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "gc" {
			return
		}
	}
	t.Fatal("gc command not registered")
}

func TestGCCmd_DefaultFlags(t *testing.T) {
	f := gcCmd.Flags()
	if v, _ := f.GetString("snapshot"); v != "snapshot.json" {
		t.Fatalf("unexpected default snapshot: %s", v)
	}
	if v, _ := f.GetBool("dry-run"); v != false {
		t.Fatal("expected dry-run default false")
	}
}
