package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSealStore(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "seals.json")
}

func TestSeal_EmptyPath(t *testing.T) {
	err := Seal("", "abc", "user", "reason")
	if err == nil || err.Error() != "seal: path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestSeal_EmptyChecksum(t *testing.T) {
	err := Seal("/tmp/x.json", "", "user", "reason")
	if err == nil || err.Error() != "seal: checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestSeal_EmptySealedBy(t *testing.T) {
	err := Seal("/tmp/x.json", "abc", "", "reason")
	if err == nil || err.Error() != "seal: sealedBy is required" {
		t.Fatalf("expected sealedBy error, got %v", err)
	}
}

func TestSeal_EmptyReason(t *testing.T) {
	err := Seal("/tmp/x.json", "abc", "user", "")
	if err == nil || err.Error() != "seal: reason is required" {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestSeal_Success(t *testing.T) {
	path := writeSealStore(t)
	err := Seal(path, "checksum1", "alice", "compliance")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sealed, err := IsSealed(path, "checksum1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sealed {
		t.Fatal("expected entry to be sealed")
	}
}

func TestSeal_DuplicatePrevented(t *testing.T) {
	path := writeSealStore(t)
	_ = Seal(path, "checksum1", "alice", "reason")
	err := Seal(path, "checksum1", "bob", "another reason")
	if err == nil || err.Error() != "seal: entry already sealed" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestIsSealed_NonExistentFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	sealed, err := IsSealed(path, "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sealed {
		t.Fatal("expected not sealed for missing file")
	}
}

func TestGetSeal_NotFound(t *testing.T) {
	path := writeSealStore(t)
	_, err := GetSeal(path, "ghost")
	if err == nil || err.Error() != "seal: record not found" {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestGetSeal_Found(t *testing.T) {
	path := writeSealStore(t)
	_ = Seal(path, "cs42", "dave", "audit")
	rec, err := GetSeal(path, "cs42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.SealedBy != "dave" || rec.Reason != "audit" {
		t.Fatalf("unexpected record: %+v", rec)
	}
}

func TestSeal_FilePermissions(t *testing.T) {
	path := writeSealStore(t)
	_ = Seal(path, "cs1", "user", "test")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
