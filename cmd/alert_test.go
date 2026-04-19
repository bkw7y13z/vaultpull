package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"vaultpull/internal/snapshot"
)

func writeAlertSnap(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "snap.json")
	entry := snapshot.Entry{
		Checksum:  "abc123",
		Keys:      []string{"KEY"},
		CreatedAt: time.Now().UTC(),
	}
	data, _ := json.Marshal(snapshot.Snapshot{Entries: []snapshot.Entry{entry}})
	os.WriteFile(path, data, 0644)
	return path
}

func TestRunAlertAdd_MissingChecksum(t *testing.T) {
	cmd := alertAddCmd
	cmd.Flags().Set("snapshot", "/tmp/snap.json")
	err := runAlertAdd(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestRunAlertAdd_MissingMessage(t *testing.T) {
	dir := t.TempDir()
	path := writeAlertSnap(t, dir)
	cmd := alertAddCmd
	cmd.Flags().Set("snapshot", path)
	cmd.Flags().Set("checksum", "abc123")
	err := runAlertAdd(cmd, nil)
	if err == nil || err.Error() != "--message is required" {
		t.Fatalf("expected message error, got %v", err)
	}
}

func TestRunAlertAdd_Success(t *testing.T) {
	dir := t.TempDir()
	path := writeAlertSnap(t, dir)
	cmd := alertAddCmd
	cmd.Flags().Set("snapshot", path)
	cmd.Flags().Set("checksum", "abc123")
	cmd.Flags().Set("message", "test alert")
	cmd.Flags().Set("by", "admin")
	cmd.Flags().Set("severity", "warning")
	err := runAlertAdd(cmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAlertList_MissingChecksum(t *testing.T) {
	cmd := alertListCmd
	cmd.Flags().Set("snapshot", "/tmp/snap.json")
	err := runAlertList(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestAlertCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "alert" {
			return
		}
	}
	t.Fatal("alert command not registered on root")
}
