package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vaultpull/internal/snapshot"
)

func writeImportSnap(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(p, nil); err != nil {
		t.Fatal(err)
	}
	return p
}

func writeSrcFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "src*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString(content)
	return f.Name()
}

func TestRunImport_MissingSource(t *testing.T) {
	_, err := executeCommand(rootCmd, "import", "--snapshot", "x.json")
	if err == nil || !strings.Contains(err.Error(), "source") {
		t.Fatalf("expected source error, got %v", err)
	}
}

func TestRunImport_Success_Env(t *testing.T) {
	snap := writeImportSnap(t)
	src := writeSrcFile(t, "TOKEN=abc\nSECRET=xyz\n")

	out, err := executeCommand(rootCmd, "import", "--snapshot", snap, "--source", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Imported") {
		t.Errorf("expected success message, got: %s", out)
	}

	entries, _ := snapshot.Load(snap)
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestRunImport_Success_JSON(t *testing.T) {
	snap := writeImportSnap(t)
	secrets := map[string]string{"K": "v"}
	b, _ := json.Marshal(secrets)
	f, _ := os.CreateTemp(t.TempDir(), "src*.json")
	f.Write(b)
	f.Close()

	_, err := executeCommand(rootCmd, "import", "--snapshot", snap, "--source", f.Name(), "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestImportCmd_RegisteredOnRoot(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "import" {
			return
		}
	}
	t.Error("import command not registered")
}

func TestImportCmd_DefaultFlags(t *testing.T) {
	f := importCmd.Flags().Lookup("snapshot")
	if f == nil || f.DefValue != "snapshot.json" {
		t.Errorf("unexpected default for --snapshot: %v", f)
	}
	ff := importCmd.Flags().Lookup("format")
	if ff == nil || ff.DefValue != "env" {
		t.Errorf("unexpected default for --format: %v", ff)
	}
}
