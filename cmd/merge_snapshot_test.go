package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func writeMergeSnap(t *testing.T, dir, name string, s snapshot.Snapshot) string {
	t.Helper()
	p := filepath.Join(dir, name)
	b, _ := json.Marshal(s)
	_ = os.WriteFile(p, b, 0600)
	return p
}

func TestRunMergeSnapshot_MissingDst(t *testing.T) {
	rootCmd.SetArgs([]string{"merge-snapshot", "--src", "x.json"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --dst")
	}
}

func TestRunMergeSnapshot_MissingSrc(t *testing.T) {
	rootCmd.SetArgs([]string{"merge-snapshot", "--dst", "x.json"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --src")
	}
}

func TestRunMergeSnapshot_Success(t *testing.T) {
	dir := t.TempDir()
	dst := writeMergeSnap(t, dir, "dst.json", snapshot.Snapshot{})
	src := writeMergeSnap(t, dir, "src.json", snapshot.Snapshot{
		Entries: []snapshot.Entry{
			{Checksum: "aaa", Keys: []string{"X"}, CreatedAt: time.Now()},
		},
	})

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"merge-snapshot", "--dst", dst, "--src", src})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("added=1")) {
		t.Errorf("expected added=1 in output, got: %s", buf.String())
	}
}
