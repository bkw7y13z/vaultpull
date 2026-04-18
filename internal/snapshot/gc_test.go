package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeGCSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "gc-*.json")
	if err != nil {
		t.Fatal(err)
	}
	snap := Snapshot{Entries: entries}
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestGC_EmptyPath(t *testing.T) {
	_, err := GC(GCOptions{MaxAge: time.Hour})
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestGC_ZeroMaxAge(t *testing.T) {
	_, err := GC(GCOptions{SnapshotPath: "x.json", MaxAge: 0})
	if err == nil {
		t.Fatal("expected error for zero max age")
	}
}

func TestGC_RemovesOldEntries(t *testing.T) {
	old := Entry{Checksum: "aaa", CreatedAt: time.Now().Add(-48 * time.Hour)}
	new := Entry{Checksum: "bbb", CreatedAt: time.Now()}
	path := writeGCSnapshot(t, []Entry{old, new})

	res, err := GC(GCOptions{SnapshotPath: path, MaxAge: 24 * time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "aaa" {
		t.Fatalf("expected aaa removed, got %v", res.Removed)
	}
}

func TestGC_KeepsPinned(t *testing.T) {
	old := Entry{Checksum: "aaa", CreatedAt: time.Now().Add(-48 * time.Hour), Pinned: true}
	path := writeGCSnapshot(t, []Entry{old})

	res, err := GC(GCOptions{SnapshotPath: path, MaxAge: time.Hour, KeepPinned: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 0 {
		t.Fatalf("expected pinned entry kept, got removed: %v", res.Removed)
	}
}

func TestGC_KeepsTagged(t *testing.T) {
	old := Entry{Checksum: "bbb", CreatedAt: time.Now().Add(-48 * time.Hour), Tag: "release"}
	path := writeGCSnapshot(t, []Entry{old})

	res, err := GC(GCOptions{SnapshotPath: path, MaxAge: time.Hour, KeepTagged: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 0 {
		t.Fatalf("expected tagged entry kept")
	}
}

func TestGC_DryRun(t *testing.T) {
	old := Entry{Checksum: "ccc", CreatedAt: time.Now().Add(-48 * time.Hour)}
	path := writeGCSnapshot(t, []Entry{old})

	res, err := GC(GCOptions{SnapshotPath: path, MaxAge: time.Hour, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 1 {
		t.Fatal("expected removal reported in dry run")
	}
	snap, _ := Load(path)
	if len(snap.Entries) != 1 {
		t.Fatal("dry run should not modify file")
	}
}
