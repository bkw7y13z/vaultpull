package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"vaultpull/internal/snapshot"
)

func writeRestoreSnap(t *testing.T, entries []snapshot.Entry) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "snap.json")
	data, _ := json.Marshal(entries)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func TestRunRestore_MissingRef(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"restore"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing ref argument")
	}
}

func TestRunRestore_InvalidSnapshot(t *testing.T) {
	output := filepath.Join(t.TempDir(), "out.env")
	cmd := rootCmd
	cmd.SetArgs([]string{"restore", "abc",
		"--snapshot", "/nonexistent/snap.json",
		"--output", output,
	})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRunRestore_Success(t *testing.T) {
	secrets := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	snap := writeRestoreSnap(t, []snapshot.Entry{
		{
			Checksum:  "abc123def456",
			Tag:       "release",
			Timestamp: time.Now(),
			Secrets:   secrets,
		},
	})
	output := filepath.Join(t.TempDir(), "out.env")

	cmd := rootCmd
	cmd.SetArgs([]string{"restore", "release",
		"--snapshot", snap,
		"--output", output,
		"--overwrite",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Error("expected non-empty .env file")
	}
}

func TestRunRestore_RefNotFound(t *testing.T) {
	snap := writeRestoreSnap(t, []snapshot.Entry{
		{Checksum: "aaa", Timestamp: time.Now(), Secrets: map[string]string{"X": "1"}},
	})
	output := filepath.Join(t.TempDir(), "out.env")
	cmd := rootCmd
	cmd.SetArgs([]string{"restore", "zzz",
		"--snapshot", snap,
		"--output", output,
	})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for unknown ref")
	}
}
