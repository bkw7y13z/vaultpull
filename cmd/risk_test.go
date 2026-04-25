package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeRiskSnap(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snapshot.json")
	entry := snapshot.Entry{
		Checksum:  "deadbeef",
		Keys:      []string{"SECRET"},
		CreatedAt: time.Now().UTC(),
	}
	data, _ := json.Marshal(map[string]interface{}{"entries": []snapshot.Entry{entry}})
	os.WriteFile(p, data, 0644)
	return p
}

func TestRunRiskAssess_MissingChecksum(t *testing.T) {
	cmd := getRiskAssessCmd()
	cmd.Flags().Set("snapshot", "/tmp/snap.json")
	cmd.Flags().Set("level", "high")
	cmd.Flags().Set("reason", "reason")
	cmd.Flags().Set("by", "user")
	err := cmd.RunE(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestRunRiskAssess_MissingLevel(t *testing.T) {
	cmd := getRiskAssessCmd()
	cmd.Flags().Set("checksum", "abc")
	cmd.Flags().Set("reason", "reason")
	cmd.Flags().Set("by", "user")
	err := cmd.RunE(cmd, nil)
	if err == nil || err.Error() != "--level is required" {
		t.Fatalf("expected level error, got %v", err)
	}
}

func TestRunRiskAssess_Success(t *testing.T) {
	p := writeRiskSnap(t)
	cmd := getRiskAssessCmd()
	cmd.Flags().Set("snapshot", p)
	cmd.Flags().Set("checksum", "deadbeef")
	cmd.Flags().Set("level", "critical")
	cmd.Flags().Set("reason", "leaked")
	cmd.Flags().Set("by", "alice")
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRiskGet_MissingChecksum(t *testing.T) {
	cmd := getRiskGetCmd()
	cmd.Flags().Set("snapshot", "/tmp/snap.json")
	err := cmd.RunE(cmd, nil)
	if err == nil || err.Error() != "--checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestRunRiskGet_Success(t *testing.T) {
	p := writeRiskSnap(t)
	snapshot.AssessRisk(p, "deadbeef", "test reason", "alice", snapshot.RiskMedium)
	cmd := getRiskGetCmd()
	cmd.Flags().Set("snapshot", p)
	cmd.Flags().Set("checksum", "deadbeef")
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func getRiskAssessCmd() *cobra.Command {
	cmd := &cobra.Command{RunE: runRiskAssess}
	cmd.Flags().String("snapshot", "", "")
	cmd.Flags().String("checksum", "", "")
	cmd.Flags().String("level", "", "")
	cmd.Flags().String("reason", "", "")
	cmd.Flags().String("by", "", "")
	return cmd
}

func getRiskGetCmd() *cobra.Command {
	cmd := &cobra.Command{RunE: runRiskGet}
	cmd.Flags().String("snapshot", "", "")
	cmd.Flags().String("checksum", "", "")
	return cmd
}
