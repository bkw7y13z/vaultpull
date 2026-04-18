package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeBlameSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "blame-*.json")
	if err != nil {
		t.Fatal(err)
	}
	store := Store{Entries: entries}
	if err := json.NewEncoder(f).Encode(store); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestBlame_EmptyPath(t *testing.T) {
	_, err := Blame("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBlame_EmptySnapshot(t *testing.T) {
	path := writeBlameSnapshot(t, nil)
	out, err := Blame(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestBlame_SingleEntry(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	path := writeBlameSnapshot(t, []Entry{
		{Checksum: "abc", Keys: []string{"DB_HOST", "DB_PORT"}, CreatedAt: now},
	})
	out, err := Blame(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 blame entries, got %d", len(out))
	}
	if out[0].Key != "DB_HOST" || out[0].Checksum != "abc" {
		t.Errorf("unexpected entry: %+v", out[0])
	}
}

func TestBlame_LatestWins(t *testing.T) {
	t1 := time.Now().Add(-time.Hour).UTC()
	t2 := time.Now().UTC()
	path := writeBlameSnapshot(t, []Entry{
		{Checksum: "old", Keys: []string{"SECRET"}, CreatedAt: t1},
		{Checksum: "new", Keys: []string{"SECRET"}, CreatedAt: t2},
	})
	out, err := Blame(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Checksum != "new" {
		t.Errorf("expected latest checksum, got %s", out[0].Checksum)
	}
}

func TestBlame_SortedKeys(t *testing.T) {
	path := writeBlameSnapshot(t, []Entry{
		{Checksum: "x", Keys: []string{"ZEBRA", "ALPHA", "MIDDLE"}, CreatedAt: time.Now()},
	})
	out, err := Blame(path)
	if err != nil {
		t.Fatal(err)
	}
	if out[0].Key != "ALPHA" || out[1].Key != "MIDDLE" || out[2].Key != "ZEBRA" {
		t.Errorf("keys not sorted: %v", out)
	}
}
