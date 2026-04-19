package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writePromoteSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	snap := Snapshot{Entries: entries}
	if err := Save(p, snap); err != nil {
		t.Fatalf("save: %v", err)
	}
	return p
}

func TestPromote_EmptyPath(t *testing.T) {
	err := Promote("", "abc", "dev", "prod", "alice", "")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestPromote_EmptyChecksum(t *testing.T) {
	err := Promote("snap.json", "", "dev", "prod", "alice", "")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestPromote_EmptyFromEnv(t *testing.T) {
	err := Promote("snap.json", "abc", "", "prod", "alice", "")
	if err == nil || err.Error() != "from_env is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestPromote_ChecksumNotFound(t *testing.T) {
	p := writePromoteSnapshot(t, []Entry{{Checksum: "aaa", Keys: []string{"X"}, CreatedAt: time.Now()}})
	err := Promote(p, "zzz", "dev", "prod", "alice", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPromote_Success(t *testing.T) {
	p := writePromoteSnapshot(t, []Entry{{Checksum: "abc123", Keys: []string{"KEY"}, CreatedAt: time.Now()}})
	if err := Promote(p, "abc123", "dev", "prod", "alice", "first promo"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	records, err := GetPromotions(p, "abc123")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	r := records[0]
	if r.FromEnv != "dev" || r.ToEnv != "prod" || r.PromotedBy != "alice" || r.Note != "first promo" {
		t.Errorf("unexpected record: %+v", r)
	}
}

func TestPromote_Appends(t *testing.T) {
	p := writePromoteSnapshot(t, []Entry{{Checksum: "abc123", Keys: []string{"KEY"}, CreatedAt: time.Now()}})
	_ = Promote(p, "abc123", "dev", "staging", "alice", "")
	_ = Promote(p, "abc123", "staging", "prod", "bob", "")
	records, _ := GetPromotions(p, "abc123")
	if len(records) != 2 {
		t.Fatalf("expected 2, got %d", len(records))
	}
}

func TestGetPromotions_EmptyPath(t *testing.T) {
	_, err := GetPromotions("", "abc")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetPromotions_NonExistentStore(t *testing.T) {
	p := writePromoteSnapshot(t, []Entry{{Checksum: "abc", Keys: []string{"K"}, CreatedAt: time.Now()}})
	// no promotions file yet
	records, err := GetPromotions(p, "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected 0, got %d", len(records))
	}
	_ = os.Remove(promotePath(p))
}
