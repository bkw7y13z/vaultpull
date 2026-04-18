package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeExpireSnapshot(t *testing.T, path string) {
	t.Helper()
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"FOO"}, CreatedAt: time.Now()},
		},
	}
	data, _ := json.Marshal(snap)
	os.WriteFile(path, data, 0600)
}

func TestSetExpiry_EmptyPath(t *testing.T) {
	err := SetExpiry("", "abc", "user", time.Now().Add(time.Hour))
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestSetExpiry_EmptyChecksum(t *testing.T) {
	err := SetExpiry("/tmp/snap", "", "user", time.Now().Add(time.Hour))
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestSetExpiry_EmptySetBy(t *testing.T) {
	err := SetExpiry("/tmp/snap", "abc", "", time.Now().Add(time.Hour))
	if err == nil || err.Error() != "set_by is required" {
		t.Fatalf("expected set_by error, got %v", err)
	}
}

func TestSetExpiry_ZeroTime(t *testing.T) {
	err := SetExpiry("/tmp/snap", "abc", "user", time.Time{})
	if err == nil || err.Error() != "expires_at is required" {
		t.Fatalf("expected time error, got %v", err)
	}
}

func TestSetExpiry_ChecksumNotFound(t *testing.T) {
	path := t.TempDir() + "/snap.json"
	writeExpireSnapshot(t, path)
	err := SetExpiry(path, "notexist", "user", time.Now().Add(time.Hour))
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestSetAndGetExpiry_Success(t *testing.T) {
	path := t.TempDir() + "/snap.json"
	writeExpireSnapshot(t, path)
	expiry := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
	if err := SetExpiry(path, "abc123", "admin", expiry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, err := GetExpiry(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected expiry entry, got nil")
	}
	if !e.ExpiresAt.Equal(expiry) {
		t.Errorf("expected %v, got %v", expiry, e.ExpiresAt)
	}
	if e.SetBy != "admin" {
		t.Errorf("expected admin, got %s", e.SetBy)
	}
}

func TestGetExpiry_NotSet(t *testing.T) {
	path := t.TempDir() + "/snap.json"
	writeExpireSnapshot(t, path)
	e, err := GetExpiry(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e != nil {
		t.Fatal("expected nil for unset expiry")
	}
}
