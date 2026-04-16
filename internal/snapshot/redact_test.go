package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeRedactSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "snap.json")
	data, _ := json.Marshal(&Store{Entries: entries})
	os.WriteFile(p, data, 0600)
	return p
}

func TestRedact_EmptyPath(t *testing.T) {
	_, err := Redact("", RedactOptions{KeyPatterns: []string{"pass"}})
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestRedact_NoPatterns(t *testing.T) {
	p := writeRedactSnapshot(t, []Entry{{Checksum: "abc", CreatedAt: time.Now()}})
	_, err := Redact(p, RedactOptions{})
	if err == nil {
		t.Fatal("expected error for missing patterns")
	}
}

func TestRedact_EmptySnapshot(t *testing.T) {
	p := writeRedactSnapshot(t, []Entry{})
	_, err := Redact(p, RedactOptions{KeyPatterns: []string{"pass"}})
	if err == nil {
		t.Fatal("expected error for empty snapshot")
	}
}

func TestRedact_MatchingKeys(t *testing.T) {
	entry := Entry{
		Checksum:  "abc123",
		CreatedAt: time.Now(),
		Secrets:   map[string]string{"DB_PASSWORD": "s3cr3t", "API_KEY": "key123", "APP_NAME": "myapp"},
	}
	p := writeRedactSnapshot(t, []Entry{entry})

	res, err := Redact(p, RedactOptions{KeyPatterns: []string{"password", "key"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.RedactedKeys) != 2 {
		t.Fatalf("expected 2 redacted keys, got %d", len(res.RedactedKeys))
	}

	store, _ := Load(p)
	latest := store.Entries[len(store.Entries)-1]
	if latest.Secrets["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected DB_PASSWORD redacted, got %q", latest.Secrets["DB_PASSWORD"])
	}
	if latest.Secrets["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unchanged")
	}
}

func TestRedact_CustomReplacement(t *testing.T) {
	entry := Entry{
		Checksum:  "xyz",
		CreatedAt: time.Now(),
		Secrets:   map[string]string{"SECRET_TOKEN": "tok"},
	}
	p := writeRedactSnapshot(t, []Entry{entry})

	_, err := Redact(p, RedactOptions{KeyPatterns: []string{"secret"}, Replacement: "[hidden]"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	store, _ := Load(p)
	if store.Entries[0].Secrets["SECRET_TOKEN"] != "[hidden]" {
		t.Errorf("expected custom replacement")
	}
}
