package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeSignSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	snap := &Snapshot{Entries: entries}
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestSignEntry_EmptyPath(t *testing.T) {
	err := SignEntry("", "abc", "secret")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestSignEntry_EmptyChecksum(t *testing.T) {
	err := SignEntry("/tmp/x.json", "", "secret")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestSignEntry_EmptySecret(t *testing.T) {
	err := SignEntry("/tmp/x.json", "abc", "")
	if err == nil || err.Error() != "signing secret is required" {
		t.Fatalf("expected secret error, got %v", err)
	}
}

func TestSignEntry_ChecksumNotFound(t *testing.T) {
	path := writeSignSnapshot(t, []Entry{{Checksum: "aaa", At: time.Now()}})
	err := SignEntry(path, "zzz", "secret")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestSignEntry_Success(t *testing.T) {
	path := writeSignSnapshot(t, []Entry{{Checksum: "abc123", At: time.Now()}})
	if err := SignEntry(path, "abc123", "mysecret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if snap.Entries[0].Metadata["signature"] == "" {
		t.Fatal("expected signature to be set")
	}
}

func TestVerifySignature_Valid(t *testing.T) {
	path := writeSignSnapshot(t, []Entry{{Checksum: "abc123", At: time.Now()}})
	if err := SignEntry(path, "abc123", "mysecret"); err != nil {
		t.Fatal(err)
	}

	ok, err := VerifySignature(path, "abc123", "mysecret")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected signature to be valid")
	}
}

func TestVerifySignature_WrongSecret(t *testing.T) {
	path := writeSignSnapshot(t, []Entry{{Checksum: "abc123", At: time.Now()}})
	if err := SignEntry(path, "abc123", "mysecret"); err != nil {
		t.Fatal(err)
	}

	ok, err := VerifySignature(path, "abc123", "wrongsecret")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected signature to be invalid")
	}
}

func TestVerifySignature_NoSignature(t *testing.T) {
	path := writeSignSnapshot(t, []Entry{{Checksum: "abc123", At: time.Now()}})
	ok, err := VerifySignature(path, "abc123", "mysecret")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected false for unsigned entry")
	}
}
