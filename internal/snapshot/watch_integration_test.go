package snapshot_test

import (
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"vaultpull/internal/snapshot"
)

func TestWatch_Integration_ChangeTriggersSnapshotDiff(t *testing.T) {
	dir := t.TempDir()
	snapPath := filepath.Join(dir, "snap.json")

	call := 0
	var changeResults []snapshot.WatchResult

	secrets := func() (map[string]string, error) {
		call++
		switch call {
		case 1:
			return map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}, nil
		case 2:
			return map[string]string{"DB_HOST": "remotehost", "DB_PORT": "5432", "DB_NAME": "prod"}, nil
		default:
			return map[string]string{"DB_HOST": "remotehost", "DB_PORT": "5432", "DB_NAME": "prod"}, nil
		}
	}

	// Pre-seed snapshot so diff has a baseline on second cycle.
	initialSecrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	snap := &snapshot.Snapshot{}
	snap.Add(snapshot.Entry{
		Timestamp: time.Now().UTC(),
		Checksum:  snapshot.ComputeChecksum(initialSecrets),
		Keys:      snapshot.KeysFromSecrets(initialSecrets),
	})
	if err := snapshot.Save(snapPath, snap); err != nil {
		t.Fatal(err)
	}

	var noChangeCalled atomic.Int32

	err := snapshot.Watch(secrets, snapshot.WatchOptions{
		SnapshotPath: snapPath,
		Interval:     time.Millisecond,
		MaxCycles:    3,
		OnChange: func(r snapshot.WatchResult) {
			changeResults = append(changeResults, r)
		},
		OnNoChange: func(_ snapshot.WatchResult) {
			noChangeCalled.Add(1)
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changeResults) == 0 {
		t.Fatal("expected at least one change result")
	}
	if noChangeCalled.Load() == 0 {
		t.Fatal("expected at least one no-change cycle")
	}
}
