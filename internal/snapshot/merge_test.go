package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeMergeSnapshot(t *testing.T, dir, name string, s Snapshot) string {
	t.Helper()
	p := filepath.Join(dir, name)
	b, _ := json.Marshal(s)
	_ = os.WriteFile(p, b, 0600)
	return p
}

func TestMerge_EmptyDstPath(t *testing.T) {
	_, err := Merge("", "src.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMerge_EmptySrcPath(t *testing.T) {
	_, err := Merge("dst.json", "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMerge_NonExistentSrc(t *testing.T) {
	dir := t.TempDir()
	dst := writeMergeSnapshot(t, dir, "dst.json", Snapshot{})
	_, err := Merge(dst, filepath.Join(dir, "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing src")
	}
}

func TestMerge_AddsNewEntries(t *testing.T) {
	dir := t.TempDir()
	dst := writeMergeSnapshot(t, dir, "dst.json", Snapshot{})
	src := writeMergeSnapshot(t, dir, "src.json", Snapshot{
		Entries: []Entry{
			{Checksum: "abc", Keys: []string{"A"}, CreatedAt: time.Now()},
			{Checksum: "def", Keys: []string{"B"}, CreatedAt: time.Now()},
		},
	})

	res, err := Merge(dst, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 2 || res.Skipped != 0 || res.Total != 2 {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestMerge_SkipsDuplicates(t *testing.T) {
	dir := t.TempDir()
	dst := writeMergeSnapshot(t, dir, "dst.json", Snapshot{
		Entries: []Entry{{Checksum: "abc", Keys: []string{"A"}, CreatedAt: time.Now()}},
	})
	src := writeMergeSnapshot(t, dir, "src.json", Snapshot{
		Entries: []Entry{
			{Checksum: "abc", Keys: []string{"A"}, CreatedAt: time.Now()},
			{Checksum: "xyz", Keys: []string{"C"}, CreatedAt: time.Now()},
		},
	})

	res, err := Merge(dst, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 1 || res.Skipped != 1 {
		t.Errorf("unexpected result: %+v", res)
	}
}
