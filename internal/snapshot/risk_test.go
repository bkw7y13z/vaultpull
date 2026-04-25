package snapshot_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeRiskSnapshot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snapshot.json")
	entry := snapshot.Entry{
		Checksum:  "abc123",
		Keys:      []string{"KEY"},
		CreatedAt: time.Now().UTC(),
	}
	data, _ := json.Marshal(map[string]interface{}{"entries": []snapshot.Entry{entry}})
	os.WriteFile(p, data, 0644)
	return p
}

func TestAssessRisk_EmptyPath(t *testing.T) {
	err := snapshot.AssessRisk("", "abc", "reason", "user", snapshot.RiskLow)
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestAssessRisk_EmptyChecksum(t *testing.T) {
	p := writeRiskSnapshot(t)
	err := snapshot.AssessRisk(p, "", "reason", "user", snapshot.RiskLow)
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestAssessRisk_EmptyReason(t *testing.T) {
	p := writeRiskSnapshot(t)
	err := snapshot.AssessRisk(p, "abc123", "", "user", snapshot.RiskLow)
	if err == nil || err.Error() != "reason is required" {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestAssessRisk_EmptyAssessedBy(t *testing.T) {
	p := writeRiskSnapshot(t)
	err := snapshot.AssessRisk(p, "abc123", "reason", "", snapshot.RiskLow)
	if err == nil || err.Error() != "assessed_by is required" {
		t.Fatalf("expected assessed_by error, got %v", err)
	}
}

func TestAssessRisk_InvalidLevel(t *testing.T) {
	p := writeRiskSnapshot(t)
	err := snapshot.AssessRisk(p, "abc123", "reason", "user", snapshot.RiskLevel("extreme"))
	if err == nil {
		t.Fatal("expected invalid risk level error")
	}
}

func TestAssessRisk_ChecksumNotFound(t *testing.T) {
	p := writeRiskSnapshot(t)
	err := snapshot.AssessRisk(p, "notfound", "reason", "user", snapshot.RiskHigh)
	if err == nil {
		t.Fatal("expected checksum not found error")
	}
}

func TestAssessRisk_Success(t *testing.T) {
	p := writeRiskSnapshot(t)
	err := snapshot.AssessRisk(p, "abc123", "exposed key", "alice", snapshot.RiskHigh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry, found, err := snapshot.GetRisk(p, "abc123")
	if err != nil {
		t.Fatalf("get risk error: %v", err)
	}
	if !found {
		t.Fatal("expected risk entry to be found")
	}
	if entry.Level != snapshot.RiskHigh {
		t.Errorf("expected high, got %s", entry.Level)
	}
	if entry.AssessedBy != "alice" {
		t.Errorf("expected alice, got %s", entry.AssessedBy)
	}
}

func TestAssessRisk_UpdatesExisting(t *testing.T) {
	p := writeRiskSnapshot(t)
	snapshot.AssessRisk(p, "abc123", "initial", "alice", snapshot.RiskLow)
	snapshot.AssessRisk(p, "abc123", "updated", "bob", snapshot.RiskCritical)

	entry, found, err := snapshot.GetRisk(p, "abc123")
	if err != nil || !found {
		t.Fatalf("expected entry, err=%v found=%v", err, found)
	}
	if entry.Level != snapshot.RiskCritical {
		t.Errorf("expected critical, got %s", entry.Level)
	}
	if entry.AssessedBy != "bob" {
		t.Errorf("expected bob, got %s", entry.AssessedBy)
	}
}

func TestGetRisk_NotFound(t *testing.T) {
	p := writeRiskSnapshot(t)
	_, found, err := snapshot.GetRisk(p, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected not found")
	}
}
