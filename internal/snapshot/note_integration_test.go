package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNote_Integration_AddAndRetrieve(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")

	snap := &Snapshot{
		Entries: []Entry{
			{Checksum: "integ001", CreatedAt: time.Now(), Keys: []string{"DB_PASS", "API_KEY"}},
			{Checksum: "integ002", CreatedAt: time.Now().Add(-time.Hour), Keys: []string{"SECRET"}},
		},
	}
	data, err := encodeSnapshot(snap)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err := os.WriteFile(p, data, 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := AddNote(p, "integ001", "production secrets"); err != nil {
		t.Fatalf("AddNote: %v", err)
	}
	if err := AddNote(p, "integ002", "staging secrets"); err != nil {
		t.Fatalf("AddNote: %v", err)
	}

	n1, err := GetNote(p, "integ001")
	if err != nil || n1 != "production secrets" {
		t.Errorf("integ001: got %q, err %v", n1, err)
	}
	n2, err := GetNote(p, "integ002")
	if err != nil || n2 != "staging secrets" {
		t.Errorf("integ002: got %q, err %v", n2, err)
	}

	// overwrite
	if err := AddNote(p, "integ001", "updated note"); err != nil {
		t.Fatalf("overwrite AddNote: %v", err)
	}
	n1b, _ := GetNote(p, "integ001")
	if n1b != "updated note" {
		t.Errorf("expected updated note, got %q", n1b)
	}
}
