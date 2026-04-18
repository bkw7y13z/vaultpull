package snapshot_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func TestLabel_Integration_SetAndRetrieve(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	snap := &snapshot.Snapshot{
		Entries: []snapshot.Entry{
			{Checksum: "chk1", CreatedAt: time.Now(), Keys: []string{"API_KEY"}},
			{Checksum: "chk2", CreatedAt: time.Now(), Keys: []string{"DB_URL"}},
		},
	}
	if err := snapshot.Save(path, snap); err != nil {
		t.Fatal(err)
	}

	if err := snapshot.Label(path, "chk1", "v1.0-release"); err != nil {
		t.Fatalf("label chk1: %v", err)
	}
	if err := snapshot.Label(path, "chk2", "staging-deploy"); err != nil {
		t.Fatalf("label chk2: %v", err)
	}

	lbl1, err := snapshot.GetLabel(path, "chk1")
	if err != nil {
		t.Fatal(err)
	}
	if lbl1 != "v1.0-release" {
		t.Errorf("expected v1.0-release, got %q", lbl1)
	}

	lbl2, err := snapshot.GetLabel(path, "chk2")
	if err != nil {
		t.Fatal(err)
	}
	if lbl2 != "staging-deploy" {
		t.Errorf("expected staging-deploy, got %q", lbl2)
	}

	// Overwrite label
	if err := snapshot.Label(path, "chk1", "v1.1-hotfix"); err != nil {
		t.Fatal(err)
	}
	updated, _ := snapshot.GetLabel(path, "chk1")
	if updated != "v1.1-hotfix" {
		t.Errorf("expected v1.1-hotfix after overwrite, got %q", updated)
	}
}
