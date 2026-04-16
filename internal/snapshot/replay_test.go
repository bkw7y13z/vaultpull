package snapshot

import (
	"path/filepath"
	"testing"
	"time"
)

func writeReplaySnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "snap.json")
	store := &Store{Entries: entries}
	if err := saveStore(p, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	return p
}

func makeEntry(checksum, tag string, secrets map[string]string) Entry {
	return Entry{
		Checksum:  checksum,
		Tag:       tag,
		CreatedAt: time.Now(),
		Secrets:   secrets,
	}
}

func TestReplay_EmptyPath(t *testing.T) {
	err := Replay("", ReplayOptions{From: "abc"}, func(ReplayEvent) error { return nil })
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestReplay_EmptyFrom(t *testing.T) {
	err := Replay("/tmp/x.json", ReplayOptions{}, func(ReplayEvent) error { return nil })
	if err == nil || err.Error() != "from ref is required" {
		t.Fatalf("expected from error, got %v", err)
	}
}

func TestReplay_NilHandler(t *testing.T) {
	err := Replay("/tmp/x.json", ReplayOptions{From: "abc"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReplay_RefNotFound(t *testing.T) {
	p := writeReplaySnapshot(t, []Entry{makeEntry("aaa", "", map[string]string{"K": "1"})})
	err := Replay(p, ReplayOptions{From: "zzz"}, func(ReplayEvent) error { return nil })
	if err == nil {
		t.Fatal("expected error for missing ref")
	}
}

func TestReplay_FullRange(t *testing.T) {
	entries := []Entry{
		makeEntry("aaa", "v1", map[string]string{"A": "1"}),
		makeEntry("bbb", "", map[string]string{"A": "1", "B": "2"}),
		makeEntry("ccc", "v3", map[string]string{"A": "1", "B": "3"}),
	}
	p := writeReplaySnapshot(t, entries)
	var events []ReplayEvent
	err := Replay(p, ReplayOptions{From: "aaa", To: "ccc"}, func(e ReplayEvent) error {
		events = append(events, e)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if events[0].Diff != nil {
		t.Error("first event should have no diff")
	}
	if events[1].Diff == nil {
		t.Error("second event should have diff")
	}
}

func TestReplay_DryRun(t *testing.T) {
	entries := []Entry{
		makeEntry("aaa", "", map[string]string{"X": "1"}),
		makeEntry("bbb", "", map[string]string{"X": "2"}),
	}
	p := writeReplaySnapshot(t, entries)
	called := false
	err := Replay(p, ReplayOptions{From: "aaa", Dry: true}, func(ReplayEvent) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("handler should not be called in dry mode")
	}
}
