package snapshot_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func TestGC_Integration_PrunesAndPreserves(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "gc-int-*.json")
	if err != nil {
		t.Fatal(err)
	}

	entries := []snapshot.Entry{
		{Checksum: "old1", CreatedAt: time.Now().Add(-72 * time.Hour)},
		{Checksum: "old2", CreatedAt: time.Now().Add(-48 * time.Hour), Pinned: true},
		{Checksum: "old3", CreatedAt: time.Now().Add(-48 * time.Hour), Tag: "v1"},
		{Checksum: "new1", CreatedAt: time.Now()},
	}
	if err := json.NewEncoder(f).Encode(snapshot.Snapshot{Entries: entries}); err != nil {
		t.Fatal(err)
	}
	f.Close()

	res, err := snapshot.GC(snapshot.GCOptions{
		SnapshotPath: f.Name(),
		MaxAge:       24 * time.Hour,
		KeepPinned:   true,
		KeepTagged:   true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Removed) != 1 || res.Removed[0] != "old1" {
		t.Fatalf("expected only old1 removed, got %v", res.Removed)
	}
	if len(res.Kept) != 3 {
		t.Fatalf("expected 3 kept, got %d", len(res.Kept))
	}

	loaded, err := snapshot.Load(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Entries) != 3 {
		t.Fatalf("expected 3 entries persisted, got %d", len(loaded.Entries))
	}
}
