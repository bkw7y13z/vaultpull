package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeAnnotateSnap(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "snap.json")
	data, _ := json.MarshalIndent(&snapshot.Snapshot{Entries: entries}, "", "  ")
	os.WriteFile(path, data, 0600)
	return path
}

func TestRunAnnotate_MissingChecksum(t *testing.T) {
	annotateChecksum = ""
	annotateNote = "hello"
	annotateGet = false
	if err := runAnnotate(nil, nil); err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestRunAnnotate_MissingNote(t *testing.T) {
	annotateChecksum = "abc"
	annotateNote = ""
	annotateGet = false
	if err := runAnnotate(nil, nil); err == nil {
		t.Fatal("expected error for missing note")
	}
}

func TestRunAnnotate_Success(t *testing.T) {
	path := writeAnnotateSnap(t, []snapshot.Entry{
		{Checksum: "deadbeef", CapturedAt: time.Now()},
	})
	snapshotPath = path
	annotateChecksum = "deadbeef"
	annotateNote = "release v1.2"
	annotateGet = false

	if err := runAnnotate(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	note, err := snapshot.GetAnnotation(path, "deadbeef")
	if err != nil {
		t.Fatalf("get annotation: %v", err)
	}
	if note != "release v1.2" {
		t.Errorf("expected %q, got %q", "release v1.2", note)
	}
}

func TestRunAnnotate_GetFlag(t *testing.T) {
	path := writeAnnotateSnap(t, []snapshot.Entry{
		{Checksum: "cafebabe", CapturedAt: time.Now(), Note: "existing note"},
	})
	snapshotPath = path
	annotateChecksum = "cafebabe"
	annotateGet = true

	if err := runAnnotate(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
