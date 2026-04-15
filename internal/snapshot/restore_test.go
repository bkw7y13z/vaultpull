package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeRestoreSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "snap.json")
	if err := Save(path, entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path
}

func TestRestore_EmptyPath(t *testing.T) {
	_, err := Restore("", "abc")
	if err == nil || err.Error() != "snapshot path must not be empty" {
		t.Fatalf("expected empty path error, got %v", err)
	}
}

func TestRestore_EmptyRef(t *testing.T) {
	_, err := Restore("/tmp/snap.json", "")
	if err == nil || err.Error() != "ref (checksum or tag) must not be empty" {
		t.Fatalf("expected empty ref error, got %v", err)
	}
}

func TestRestore_NonExistentFile(t *testing.T) {
	_, err := Restore(filepath.Join(t.TempDir(), "missing.json"), "abc")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRestore_RefNotFound(t *testing.T) {
	path := writeRestoreSnapshot(t, []Entry{
		{Checksum: "aabbcc", Timestamp: time.Now(), Secrets: map[string]string{"K": "v"}},
	})
	_, err := Restore(path, "zzz")
	if err == nil {
		t.Fatal("expected error for unknown ref")
	}
}

func TestRestore_ByChecksum(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	path := writeRestoreSnapshot(t, []Entry{
		{Checksum: "deadbeef1234", Timestamp: time.Now(), Secrets: secrets},
	})
	res, err := Restore(path, "deadbeef1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Checksum != "deadbeef1234" {
		t.Errorf("expected checksum deadbeef1234, got %s", res.Checksum)
	}
	if len(res.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Keys))
	}
}

func TestRestore_ByTag(t *testing.T) {
	secrets := map[string]string{"API_KEY": "secret"}
	path := writeRestoreSnapshot(t, []Entry{
		{Checksum: "cafebabe", Tag: "v1.0", Timestamp: time.Now(), Secrets: secrets},
	})
	res, err := Restore(path, "v1.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Tag != "v1.0" {
		t.Errorf("expected tag v1.0, got %s", res.Tag)
	}
}

func TestRestore_ByChecksumPrefix(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	path := writeRestoreSnapshot(t, []Entry{
		{Checksum: "1234567890abcdef", Timestamp: time.Now(), Secrets: secrets},
	})
	res, err := Restore(path, "1234567")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
	os.Remove(path)
}
