package snapshot_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/fvbommel/vaultpull/internal/snapshot"
)

func writeTrendSnapshot(t *testing.T, store snapshot.Store) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "trend-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(store); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestTrend_EmptyPath(t *testing.T) {
	_, err := snapshot.Trend("", time.Time{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTrend_EmptySnapshot(t *testing.T) {
	path := writeTrendSnapshot(t, snapshot.Store{})
	points, err := snapshot.Trend(path, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 0 {
		t.Fatalf("expected 0 points, got %d", len(points))
	}
}

func TestTrend_AllEntries(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	store := snapshot.Store{
		Entries: []snapshot.Entry{
			{Checksum: "aaa", CreatedAt: now.Add(-2 * time.Hour), Keys: []string{"A", "B"}},
			{Checksum: "bbb", CreatedAt: now.Add(-1 * time.Hour), Keys: []string{"A", "B", "C"}},
			{Checksum: "ccc", CreatedAt: now, Keys: []string{"A"}},
		},
	}
	path := writeTrendSnapshot(t, store)
	points, err := snapshot.Trend(path, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(points))
	}
	if points[0].KeyCount != 2 || points[1].KeyCount != 3 || points[2].KeyCount != 1 {
		t.Fatalf("unexpected key counts: %+v", points)
	}
}

func TestTrend_SinceFilter(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	store := snapshot.Store{
		Entries: []snapshot.Entry{
			{Checksum: "aaa", CreatedAt: now.Add(-3 * time.Hour), Keys: []string{"A"}},
			{Checksum: "bbb", CreatedAt: now.Add(-1 * time.Hour), Keys: []string{"A", "B"}},
		},
	}
	path := writeTrendSnapshot(t, store)
	points, err := snapshot.Trend(path, now.Add(-2*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 1 {
		t.Fatalf("expected 1 point, got %d", len(points))
	}
	if points[0].Checksum != "bbb" {
		t.Fatalf("unexpected checksum: %s", points[0].Checksum)
	}
}
