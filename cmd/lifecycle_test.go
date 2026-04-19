package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaultpull/internal/snapshot"
)

func writeLifecycleSnap(t *testing.T, dir string) string {
	t.Helper()
	snap := snapshot.Snapshot{
		Entries: []snapshot.Entry{
			{Checksum: "abc123", Keys: []string{"KEY"}, CreatedAt: time.Now()},
		},
	}
	data, _ := json.Marshal(snap)
	p := filepath.Join(dir, "snap.json")
	os.WriteFile(p, data, 0644)
	return p
}

func TestRunLifecycleSet_MissingChecksum(t *testing.T) {
	cmd := lifecycleSetCmd
	cmd.Flags().Set("snapshot", "/tmp/x.json")
	err := runLifecycleSet(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestRunLifecycleSet_MissingState(t *testing.T) {
	cmd := lifecycleSetCmd
	cmd.Flags().Set("checksum", "abc")
	cmd.Flags().Set("state", "")
	err := runLifecycleSet(cmd, nil)
	if err == nil || err.Error() != "--state is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestRunLifecycleSet_Success(t *testing.T) {
	dir := t.TempDir()
	p := writeLifecycleSnap(t, dir)

	cmd := lifecycleSetCmd
	cmd.Flags().Set("snapshot", p)
	cmd.Flags().Set("checksum", "abc123")
	cmd.Flags().Set("state", "deprecated")
	cmd.Flags().Set("by", "tester")
	cmd.Flags().Set("reason", "test run")

	err := runLifecycleSet(cmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunLifecycleGet_MissingChecksum(t *testing.T) {
	cmd := lifecycleGetCmd
	cmd.Flags().Set("checksum", "")
	err := runLifecycleGet(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestLifecycleCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "lifecycle" {
			return
		}
	}
	t.Fatal("lifecycle command not registered on root")
}
