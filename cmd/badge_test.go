package cmd

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/vaultpull/vaultpull/internal/snapshot"
)

func writeBadgeSnap(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	s := snapshot.Snapshot{Entries: entries}
	if err := json.NewEncoder(f).Encode(s); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRunBadgeSet_MissingChecksum(t *testing.T) {
	cmd := badgeSetCmd
	cmd.Flags().Set("checksum", "")
	err := runBadgeSet(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestRunBadgeSet_MissingLabel(t *testing.T) {
	cmd := badgeSetCmd
	cmd.Flags().Set("checksum", "abc")
	cmd.Flags().Set("label", "")
	err := runBadgeSet(cmd, nil)
	if err == nil || err.Error() != "--label is required" {
		t.Fatalf("expected label error, got %v", err)
	}
}

func TestRunBadgeSet_Success(t *testing.T) {
	path := writeBadgeSnap(t, []snapshot.Entry{{Checksum: "deadbeef", CreatedAt: time.Now()}})
	cmd := badgeSetCmd
	cmd.Flags().Set("snapshot", path)
	cmd.Flags().Set("checksum", "deadbeef")
	cmd.Flags().Set("label", "deploy")
	cmd.Flags().Set("status", "ok")
	cmd.Flags().Set("message", "deployed successfully")
	if err := runBadgeSet(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunBadgeGet_MissingChecksum(t *testing.T) {
	cmd := badgeGetCmd
	cmd.Flags().Set("checksum", "")
	err := runBadgeGet(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestBadgeCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "badge" {
			return
		}
	}
	t.Fatal("badge command not registered on root")
}
