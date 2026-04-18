package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeLabelSnap(t *testing.T, checksum string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	snap := &snapshot.Snapshot{
		Entries: []snapshot.Entry{
			{Checksum: checksum, CreatedAt: time.Now(), Keys: []string{"DB_PASS"}},
		},
	}
	data, _ := json.MarshalIndent(snap, "", "  ")
	os.WriteFile(path, data, 0644)
	return path
}

func TestRunLabelAdd_MissingChecksum(t *testing.T) {
	args := []string{"label", "add", "--label", "prod"}
	if err := executeCmd(args); err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestRunLabelAdd_MissingLabel(t *testing.T) {
	args := []string{"label", "add", "--checksum", "abc"}
	if err := executeCmd(args); err == nil || err.Error() != "--label is required" {
		t.Fatalf("expected label error, got %v", err)
	}
}

func TestRunLabelAdd_Success(t *testing.T) {
	path := writeLabelSnap(t, "abc123")
	args := []string{"label", "add", "--snapshot", path, "--checksum", "abc123", "--label", "baseline"}
	if err := executeCmd(args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lbl, err := snapshot.GetLabel(path, "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if lbl != "baseline" {
		t.Errorf("expected 'baseline', got %q", lbl)
	}
}

func TestRunLabelGet_MissingChecksum(t *testing.T) {
	args := []string{"label", "get"}
	if err := executeCmd(args); err == nil {
		t.Fatal("expected error")
	}
}

func TestLabelCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "label" {
			found = true
			break
		}
	}
	if !found {
		t.Error("label command not registered on root")
	}
}
