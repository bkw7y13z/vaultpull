package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestQuota_Integration_SetAndRetrieve(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/snap.json"

	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "intchk1", Keys: []string{"SECRET_KEY", "API_TOKEN"}, CreatedAt: time.Now()},
			{Checksum: "intchk2", Keys: []string{"DB_PASS"}, CreatedAt: time.Now()},
		},
	}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(path, data, 0644)

	if err := SetQuota(path, "intchk1", "sre-team", 100, 256); err != nil {
		t.Fatalf("SetQuota failed: %v", err)
	}
	if err := SetQuota(path, "intchk2", "sre-team", 0, 64); err != nil {
		t.Fatalf("SetQuota failed: %v", err)
	}

	p1, ok, err := GetQuota(path, "intchk1")
	if err != nil || !ok {
		t.Fatalf("GetQuota intchk1 failed: err=%v ok=%v", err, ok)
	}
	if p1.MaxKeys != 100 || p1.MaxSizeKB != 256 {
		t.Errorf("unexpected policy for intchk1: %+v", p1)
	}

	p2, ok, err := GetQuota(path, "intchk2")
	if err != nil || !ok {
		t.Fatalf("GetQuota intchk2 failed: err=%v ok=%v", err, ok)
	}
	if p2.MaxSizeKB != 64 {
		t.Errorf("unexpected max_size_kb for intchk2: %d", p2.MaxSizeKB)
	}

	_, ok, err = GetQuota(path, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected quota not found for nonexistent checksum")
	}
}
