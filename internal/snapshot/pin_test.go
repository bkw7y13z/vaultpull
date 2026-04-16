package snapshot

import (
	"testing"
	"time"
)

func writePinSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	p := tmpPath(t)
	snap := &Snapshot{Entries: entries}
	if err := Save(p, snap); err != nil {
		t.Fatalf("save: %v", err)
	}
	return p
}

func TestPin_EmptyPath(t *testing.T) {
	if err := Pin("", "abc", "keep"); err == nil {
		t.Fatal("expected error")
	}
}

func TestPin_EmptyChecksum(t *testing.T) {
	if err := Pin("/tmp/snap.json", "", "keep"); err == nil {
		t.Fatal("expected error")
	}
}

func TestPin_EmptyReason(t *testing.T) {
	if err := Pin("/tmp/snap.json", "abc", ""); err == nil {
		t.Fatal("expected error")
	}
}

func TestPin_ChecksumNotFound(t *testing.T) {
	p := writePinSnapshot(t, []Entry{{Checksum: "aaa", CreatedAt: time.Now()}})
	if err := Pin(p, "zzz", "reason"); err == nil {
		t.Fatal("expected error")
	}
}

func TestPin_Success(t *testing.T) {
	p := writePinSnapshot(t, []Entry{{Checksum: "aaa", CreatedAt: time.Now()}})
	if err := Pin(p, "aaa", "critical release"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ok, err := IsPinned(p, "aaa")
	if err != nil {
		t.Fatalf("IsPinned: %v", err)
	}
	if !ok {
		t.Fatal("expected entry to be pinned")
	}
}

func TestUnpin_Success(t *testing.T) {
	p := writePinSnapshot(t, []Entry{{Checksum: "bbb", CreatedAt: time.Now()}})
	if err := Pin(p, "bbb", "reason"); err != nil {
		t.Fatalf("pin: %v", err)
	}
	if err := Unpin(p, "bbb"); err != nil {
		t.Fatalf("unpin: %v", err)
	}
	ok, err := IsPinned(p, "bbb")
	if err != nil {
		t.Fatalf("IsPinned: %v", err)
	}
	if ok {
		t.Fatal("expected entry to be unpinned")
	}
}

func TestIsPinned_ChecksumNotFound(t *testing.T) {
	p := writePinSnapshot(t, []Entry{{Checksum: "ccc", CreatedAt: time.Now()}})
	_, err := IsPinned(p, "zzz")
	if err == nil {
		t.Fatal("expected error")
	}
}
