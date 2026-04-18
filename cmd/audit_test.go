package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"vaultpull/internal/snapshot"
)

func writeAuditSnap(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snapshots.json")
	data, _ := json.Marshal(map[string]any{"entries": []any{}})
	os.WriteFile(p, data, 0600)
	return p
}

func TestRunAuditRecord_MissingAction(t *testing.T) {
	p := writeAuditSnap(t)
	cmd := auditRecordCmd
	cmd.Flags().Set("snapshot", p)
	cmd.Flags().Set("action", "")
	err := runAuditRecord(cmd, nil)
	if err == nil {
		t.Fatal("expected error for missing action")
	}
}

func TestRunAuditRecord_Success(t *testing.T) {
	p := writeAuditSnap(t)
	cmd := auditRecordCmd
	cmd.Flags().Set("snapshot", p)
	cmd.Flags().Set("action", "pull")
	cmd.Flags().Set("checksum", "abc")
	cmd.Flags().Set("actor", "ci")
	cmd.Flags().Set("detail", "test run")
	err := runAuditRecord(cmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events, _ := snapshot.GetAuditLog(p)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

func TestRunAuditLog_Empty(t *testing.T) {
	p := writeAuditSnap(t)
	cmd := auditLogCmd
	cmd.Flags().Set("snapshot", p)
	err := runAuditLog(cmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuditCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "audit" {
			found = true
			break
		}
	}
	if !found {
		t.Error("audit command not registered on root")
	}
}
