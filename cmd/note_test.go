package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"vaultpull/internal/snapshot"
)

func writeNoteSnap(t *testing.T, checksum string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	snap := map[string]interface{}{
		"entries": []map[string]interface{}{
			{"checksum": checksum, "created_at": time.Now(), "keys": []string{"KEY"}},
		},
	}
	data, _ := json.MarshalIndent(snap, "", "  ")
	os.WriteFile(p, data, 0600)
	return p
}

func TestRunNoteAdd_MissingChecksum(t *testing.T) {
	cmd := noteAddCmd
	cmd.Flags().Set("checksum", "")
	cmd.Flags().Set("note", "hello")
	// checksum guard handled via os.Exit; test via snapshot directly
	err := snapshot.AddNote("", "", "hello")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRunNoteAdd_MissingNote(t *testing.T) {
	err := snapshot.AddNote("", "abc", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRunNoteAdd_Success(t *testing.T) {
	p := writeNoteSnap(t, "deadbeef")
	err := snapshot.AddNote(p, "deadbeef", "my note")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	note, err := snapshot.GetNote(p, "deadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note != "my note" {
		t.Errorf("expected 'my note', got %q", note)
	}
}

func TestRunNoteGet_MissingChecksum(t *testing.T) {
	_, err := snapshot.GetNote("", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNoteCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "note" {
			return
		}
	}
	t.Fatal("note command not registered on root")
}
