package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeAccessSnap(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	data := make(map[string]snapshot.Entry)
	for _, e := range entries {
		data[e.Checksum] = e
	}
	b, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile(p, b, 0644)
	return p
}

func TestRunAccessRecord_MissingChecksum(t *testing.T) {
	p := writeAccessSnap(t, nil)
	rootCmd.SetArgs([]string{"access", "record", "--snapshot", p, "--by", "alice", "--action", "read"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestRunAccessRecord_MissingBy(t *testing.T) {
	p := writeAccessSnap(t, nil)
	rootCmd.SetArgs([]string{"access", "record", "--snapshot", p, "--checksum", "abc", "--action", "read"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --by")
	}
}

func TestRunAccessRecord_MissingAction(t *testing.T) {
	p := writeAccessSnap(t, nil)
	rootCmd.SetArgs([]string{"access", "record", "--snapshot", p, "--checksum", "abc", "--by", "alice"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --action")
	}
}

func TestRunAccessRecord_Success(t *testing.T) {
	p := writeAccessSnap(t, []snapshot.Entry{{Checksum: "abc123", CreatedAt: time.Now()}})
	rootCmd.SetArgs([]string{"access", "record",
		"--snapshot", p,
		"--checksum", "abc123",
		"--by", "alice",
		"--action", "read",
		"--reason", "ci-check",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, err := snapshot.GetAccessLog(p, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestRunAccessList_MissingChecksum(t *testing.T) {
	p := writeAccessSnap(t, nil)
	rootCmd.SetArgs([]string{"access", "list", "--snapshot", p})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestAccessCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "access" {
			return
		}
	}
	t.Fatal("access command not registered on root")
}
