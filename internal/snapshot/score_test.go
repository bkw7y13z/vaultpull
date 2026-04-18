package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeScoreSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snapshot.json")
	snap := &Snapshot{Entries: entries}
	if err := Save(path, snap); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func TestScore_EmptyPath(t *testing.T) {
	_, err := Score("", "abc")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestScore_EmptyChecksum(t *testing.T) {
	_, err := Score("/tmp/snap.json", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestScore_ChecksumNotFound(t *testing.T) {
	path := writeScoreSnapshot(t, []Entry{
		{Checksum: "aaa", Keys: []string{"FOO"}, CreatedAt: time.Now()},
	})
	_, err := Score(path, "zzz")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestScore_PerfectScore(t *testing.T) {
	path := writeScoreSnapshot(t, []Entry{
		{
			Checksum:  "abc123",
			Tag:       "v1",
			Keys:      []string{"DB_PASS"},
			CreatedAt: time.Now(),
		},
	})
	se, err := Score(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if se.Score != 100 {
		t.Errorf("expected 100, got %d (reasons: %v)", se.Score, se.Reasons)
	}
}

func TestScore_PenaltiesApplied(t *testing.T) {
	path := writeScoreSnapshot(t, []Entry{
		{
			Checksum:  "def456",
			CreatedAt: time.Time{}, // zero
		},
	})
	se, err := Score(path, "def456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if se.Score >= 100 {
		t.Errorf("expected penalties, got score %d", se.Score)
	}
	if len(se.Reasons) == 0 {
		t.Error("expected reasons")
	}
}

func TestGetScore_NotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snapshot.json")
	_, err := GetScore(path, "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetScore_AfterScore(t *testing.T) {
	path := writeScoreSnapshot(t, []Entry{
		{Checksum: "xyz", Tag: "prod", Keys: []string{"A"}, CreatedAt: time.Now()},
	})
	if _, err := Score(path, "xyz"); err != nil {
		t.Fatalf("score: %v", err)
	}
	se, err := GetScore(path, "xyz")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if se.Checksum != "xyz" {
		t.Errorf("wrong checksum: %s", se.Checksum)
	}
	// scores.json should exist
	sp := scorePath(path)
	if _, err := os.Stat(sp); err != nil {
		t.Errorf("scores file missing: %v", err)
	}
}
