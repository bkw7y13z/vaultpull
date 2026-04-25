package audit

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLogger_NoOp(t *testing.T) {
	l := NewLogger("")
	if err := l.Success("/secret/app", ".env", 5); err != nil {
		t.Fatalf("no-op logger should not return error, got: %v", err)
	}
}

func TestLogger_Success(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	l := NewLogger(path)

	if err := l.Success("/secret/app", ".env", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := readEntries(t, path)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != "success" {
		t.Errorf("expected status=success, got %s", entries[0].Status)
	}
	if entries[0].KeysCount != 3 {
		t.Errorf("expected keys_count=3, got %d", entries[0].KeysCount)
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogger_Failure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	l := NewLogger(path)

	if err := l.Failure("/secret/app", ".env", errors.New("connection refused")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := readEntries(t, path)
	if entries[0].Status != "failure" {
		t.Errorf("expected status=failure, got %s", entries[0].Status)
	}
	if entries[0].Message != "connection refused" {
		t.Errorf("expected message=connection refused, got %s", entries[0].Message)
	}
}

func TestLogger_Appends(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	l := NewLogger(path)

	_ = l.Success("/secret/a", ".env", 1)
	_ = l.Success("/secret/b", ".env", 2)

	entries := readEntries(t, path)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestLogger_SuccessRecordsPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	l := NewLogger(path)

	const secretPath = "/secret/app"
	const destFile = ".env"

	if err := l.Success(secretPath, destFile, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := readEntries(t, path)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Path != secretPath {
		t.Errorf("expected path=%s, got %s", secretPath, entries[0].Path)
	}
	if entries[0].Dest != destFile {
		t.Errorf("expected dest=%s, got %s", destFile, entries[0].Dest)
	}
}

func readEntries(t *testing.T, path string) []Entry {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("opening log file: %v", err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			t.Fatalf("unmarshalling entry: %v", err)
		}
		entries = append(entries, e)
	}
	return entries
}
