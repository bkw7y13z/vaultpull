package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeTestSnapshot(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	data, _ := json.Marshal(entries)
	_ = os.WriteFile(p, data, 0600)
	return p
}

func TestRunTag_MissingLabel(t *testing.T) {
	tagLabel = ""
	tagFind = false
	tagChecksum = "abc"
	tagSnapshotPath = "/tmp/snap.json"

	err := runTag(tagCmd, nil)
	if err == nil {
		t.Fatal("expected error for missing label")
	}
}

func TestRunTag_MissingChecksumWhenTagging(t *testing.T) {
	tagLabel = "v1"
	tagFind = false
	tagChecksum = ""
	tagSnapshotPath = "/tmp/snap.json"

	err := runTag(tagCmd, nil)
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestRunTag_Success(t *testing.T) {
	p := writeTestSnapshot(t, []snapshot.Entry{
		{Timestamp: time.Now(), Checksum: "deadbeef", Keys: []string{"API_KEY"}},
	})

	tagSnapshotPath = p
	tagChecksum = "deadbeef"
	tagLabel = "staging-v2"
	tagFind = false

	if err := runTag(tagCmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := snapshot.Load(p)
	if entries[0].Tag != "staging-v2" {
		t.Errorf("expected tag 'staging-v2', got %q", entries[0].Tag)
	}
}

func TestRunTag_FindSuccess(t *testing.T) {
	p := writeTestSnapshot(t, []snapshot.Entry{
		{Timestamp: time.Now(), Checksum: "cafebabe", Keys: []string{"DB_PASS"}, Tag: "prod-v1"},
	})

	tagSnapshotPath = p
	tagLabel = "prod-v1"
	tagFind = true

	if err := runTag(tagCmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunTag_FindNotFound(t *testing.T) {
	p := writeTestSnapshot(t, []snapshot.Entry{
		{Timestamp: time.Now(), Checksum: "cafebabe", Tag: "other"},
	})

	tagSnapshotPath = p
	tagLabel = "missing"
	tagFind = true

	if err := runTag(tagCmd, nil); err == nil {
		t.Fatal("expected error for missing tag")
	}
}
