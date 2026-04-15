package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeEnvFile_EmptyPath(t *testing.T) {
	_, err := MergeEnvFile("", map[string]string{"KEY": "val"}, false)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestMergeEnvFile_NewFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result, err := MergeEnvFile(tmp, secrets, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}
	if len(result.Updated) != 0 || len(result.Skipped) != 0 {
		t.Errorf("expected no updates or skips")
	}
}

func TestMergeEnvFile_SkipsExistingWithoutOverwrite(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(tmp, []byte("FOO=original\n"), 0600)

	result, err := MergeEnvFile(tmp, map[string]string{"FOO": "new", "BAR": "added"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "FOO" {
		t.Errorf("expected FOO to be skipped, got %v", result.Skipped)
	}
	if len(result.Added) != 1 || result.Added[0] != "BAR" {
		t.Errorf("expected BAR to be added, got %v", result.Added)
	}

	data, _ := os.ReadFile(tmp)
	if string(data) != "FOO=original\nBAR=added\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestMergeEnvFile_UpdatesExistingWithOverwrite(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(tmp, []byte("FOO=original\n"), 0600)

	result, err := MergeEnvFile(tmp, map[string]string{"FOO": "updated"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Updated) != 1 || result.Updated[0] != "FOO" {
		t.Errorf("expected FOO to be updated, got %v", result.Updated)
	}

	data, _ := os.ReadFile(tmp)
	if string(data) != "FOO=updated\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestMergeEnvFile_QuotesSpecialValues(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")

	_, err := MergeEnvFile(tmp, map[string]string{"KEY": "hello world"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	if string(data) != "KEY=\"hello world\"\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}
