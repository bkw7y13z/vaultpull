package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeAlertSnapshot(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "snap.json")
	entry := Entry{
		Checksum:  "abc123",
		Keys:      []string{"KEY"},
		CreatedAt: time.Now().UTC(),
	}
	data, _ := json.Marshal(Snapshot{Entries: []Entry{entry}})
	os.WriteFile(path, data, 0644)
	return path
}

func TestAddAlert_EmptyPath(t *testing.T) {
	err := AddAlert("", "abc123", "msg", "user", AlertWarning)
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddAlert_EmptyChecksum(t *testing.T) {
	dir := t.TempDir()
	path := writeAlertSnapshot(t, dir)
	err := AddAlert(path, "", "msg", "user", AlertWarning)
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddAlert_EmptyMessage(t *testing.T) {
	dir := t.TempDir()
	path := writeAlertSnapshot(t, dir)
	err := AddAlert(path, "abc123", "", "user", AlertWarning)
	if err == nil || err.Error() != "message is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestAddAlert_InvalidSeverity(t *testing.T) {
	dir := t.TempDir()
	path := writeAlertSnapshot(t, dir)
	err := AddAlert(path, "abc123", "msg", "user", AlertSeverity("unknown"))
	if err == nil {
		t.Fatal("expected error for invalid severity")
	}
}

func TestAddAlert_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	path := writeAlertSnapshot(t, dir)
	err := AddAlert(path, "notexist", "msg", "user", AlertInfo)
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestAddAlert_Success(t *testing.T) {
	dir := t.TempDir()
	path := writeAlertSnapshot(t, dir)
	err := AddAlert(path, "abc123", "disk usage high", "ops", AlertCritical)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	alerts, err := GetAlerts(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Severity != AlertCritical {
		t.Errorf("expected critical, got %s", alerts[0].Severity)
	}
	if alerts[0].Message != "disk usage high" {
		t.Errorf("unexpected message: %s", alerts[0].Message)
	}
}

func TestGetAlerts_EmptyPath(t *testing.T) {
	_, err := GetAlerts("", "abc123")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAlerts_NoAlerts(t *testing.T) {
	dir := t.TempDir()
	path := writeAlertSnapshot(t, dir)
	alerts, err := GetAlerts(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}
