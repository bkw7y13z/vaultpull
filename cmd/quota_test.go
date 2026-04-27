package cmd

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/username/vaultpull/internal/snapshot"
)

func writeQuotaSnap(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/snap.json"
	snap := snapshot.Snapshot{
		Entries: []snapshot.Entry{
			{Checksum: "qchk1", Keys: []string{"DB_PASS"}, CreatedAt: time.Now()},
		},
	}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestRunQuotaSet_MissingChecksum(t *testing.T) {
	cmd := quotaSetCmd
	cmd.Flags().Set("snapshot", writeQuotaSnap(t))
	cmd.Flags().Set("checksum", "")
	cmd.Flags().Set("by", "admin")
	cmd.Flags().Set("max-keys", "10")
	err := runQuotaSet(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestRunQuotaSet_MissingBy(t *testing.T) {
	path := writeQuotaSnap(t)
	cmd := quotaSetCmd
	cmd.Flags().Set("snapshot", path)
	cmd.Flags().Set("checksum", "qchk1")
	cmd.Flags().Set("by", "")
	cmd.Flags().Set("max-keys", "10")
	err := runQuotaSet(cmd, nil)
	if err == nil || err.Error() != "--by is required" {
		t.Fatalf("expected by error, got %v", err)
	}
}

func TestRunQuotaSet_Success(t *testing.T) {
	path := writeQuotaSnap(t)
	cmd := quotaSetCmd
	cmd.Flags().Set("snapshot", path)
	cmd.Flags().Set("checksum", "qchk1")
	cmd.Flags().Set("by", "admin")
	cmd.Flags().Set("max-keys", "20")
	cmd.Flags().Set("max-size-kb", "512")
	err := runQuotaSet(cmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunQuotaGet_MissingChecksum(t *testing.T) {
	path := writeQuotaSnap(t)
	cmd := quotaGetCmd
	cmd.Flags().Set("snapshot", path)
	cmd.Flags().Set("checksum", "")
	err := runQuotaGet(cmd, nil)
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestRunQuotaGet_Success(t *testing.T) {
	path := writeQuotaSnap(t)
	_ = snapshot.SetQuota(path, "qchk1", "ops", 5, 0)
	cmd := quotaGetCmd
	cmd.Flags().Set("snapshot", path)
	cmd.Flags().Set("checksum", "qchk1")
	err := runQuotaGet(cmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
