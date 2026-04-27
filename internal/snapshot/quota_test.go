package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeQuotaSnapshot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/snap.json"
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"KEY"}, CreatedAt: time.Now()},
		},
	}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestSetQuota_EmptyPath(t *testing.T) {
	err := SetQuota("", "abc123", "admin", 10, 0)
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestSetQuota_EmptyChecksum(t *testing.T) {
	path := writeQuotaSnapshot(t)
	err := SetQuota(path, "", "admin", 10, 0)
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestSetQuota_EmptySetBy(t *testing.T) {
	path := writeQuotaSnapshot(t)
	err := SetQuota(path, "abc123", "", 10, 0)
	if err == nil || err.Error() != "set_by is required" {
		t.Fatalf("expected set_by error, got %v", err)
	}
}

func TestSetQuota_NoLimits(t *testing.T) {
	path := writeQuotaSnapshot(t)
	err := SetQuota(path, "abc123", "admin", 0, 0)
	if err == nil {
		t.Fatal("expected error for zero limits")
	}
}

func TestSetQuota_ChecksumNotFound(t *testing.T) {
	path := writeQuotaSnapshot(t)
	err := SetQuota(path, "notexist", "admin", 10, 0)
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestSetQuota_Success(t *testing.T) {
	path := writeQuotaSnapshot(t)
	err := SetQuota(path, "abc123", "admin", 50, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok, err := GetQuota(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected quota to be found")
	}
	if p.MaxKeys != 50 || p.MaxSizeKB != 100 {
		t.Errorf("unexpected quota: %+v", p)
	}
	if p.SetBy != "admin" {
		t.Errorf("expected set_by=admin, got %s", p.SetBy)
	}
}

func TestGetQuota_EmptyPath(t *testing.T) {
	_, _, err := GetQuota("", "abc123")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetQuota_NotFound(t *testing.T) {
	path := writeQuotaSnapshot(t)
	_, ok, err := GetQuota(path, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected quota not found")
	}
}
