package cmd

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writePinSnap(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap*.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = json.NewEncoder(f).Encode(&snapshot.Snapshot{Entries: entries})
	f.Close()
	return f.Name()
}

func TestRunPinAdd_MissingChecksum(t *testing.T) {
	pinChecksum = ""
	pinReason = "reason"
	if err := runPinAdd(pinAddCmd, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestRunPinAdd_MissingReason(t *testing.T) {
	pinChecksum = "abc"
	pinReason = ""
	if err := runPinAdd(pinAddCmd, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestRunPinAdd_Success(t *testing.T) {
	p := writePinSnap(t, []snapshot.Entry{{Checksum: "deadbeef", CreatedAt: time.Now()}})
	pinSnapshotPath = p
	pinChecksum = "deadbeef"
	pinReason = "stable"
	if err := runPinAdd(pinAddCmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ok, err := snapshot.IsPinned(p, "deadbeef")
	if err != nil || !ok {
		t.Fatal("expected entry to be pinned")
	}
}

func TestRunPinRemove_MissingChecksum(t *testing.T) {
	pinChecksum = ""
	if err := runPinRemove(pinRemoveCmd, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestRunPinRemove_Success(t *testing.T) {
	p := writePinSnap(t, []snapshot.Entry{{Checksum: "cafebabe", CreatedAt: time.Now()}})
	pinSnapshotPath = p
	pinChecksum = "cafebabe"
	pinReason = "reason"
	_ = runPinAdd(pinAddCmd, nil)

	pinChecksum = "cafebabe"
	if err := runPinRemove(pinRemoveCmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ok, _ := snapshot.IsPinned(p, "cafebabe")
	if ok {
		t.Fatal("expected entry to be unpinned")
	}
}
