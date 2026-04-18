package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeImportSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "snap.json")
	if err := Save(p, entries); err != nil {
		t.Fatal(err)
	}
	return p
}

func writeSourceFile(t *testing.T, content, ext string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "source*"+ext)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString(content)
	return f.Name()
}

func TestImport_EmptySnapshotPath(t *testing.T) {
	err := Import(ImportOptions{SourcePath: "x"})
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestImport_EmptySourcePath(t *testing.T) {
	err := Import(ImportOptions{SnapshotPath: "x"})
	if err == nil || err.Error() != "source path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestImport_EnvFormat(t *testing.T) {
	snap := writeImportSnapshot(t, nil)
	src := writeSourceFile(t, "FOO=bar\nBAZ=qux\n", ".env")

	if err := Import(ImportOptions{SnapshotPath: snap, SourcePath: src}); err != nil {
		t.Fatal(err)
	}

	entries, _ := Load(snap)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(entries[0].Keys))
	}
}

func TestImport_JSONFormat(t *testing.T) {
	snap := writeImportSnapshot(t, nil)
	secrets := map[string]string{"KEY": "val"}
	b, _ := json.Marshal(secrets)
	src := writeSourceFile(t, string(b), ".json")

	if err := Import(ImportOptions{SnapshotPath: snap, SourcePath: src, Format: "json"}); err != nil {
		t.Fatal(err)
	}

	entries, _ := Load(snap)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestImport_DuplicateBlocked(t *testing.T) {
	snap := writeImportSnapshot(t, nil)
	src := writeSourceFile(t, "A=1\n", ".env")

	Import(ImportOptions{SnapshotPath: snap, SourcePath: src})
	err := Import(ImportOptions{SnapshotPath: snap, SourcePath: src})
	if err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestImport_OverwriteAllowed(t *testing.T) {
	snap := writeImportSnapshot(t, nil)
	src := writeSourceFile(t, "A=1\n", ".env")

	Import(ImportOptions{SnapshotPath: snap, SourcePath: src})
	err := Import(ImportOptions{SnapshotPath: snap, SourcePath: src, Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
}

func TestImport_UnsupportedFormat(t *testing.T) {
	snap := writeImportSnapshot(t, nil)
	src := writeSourceFile(t, "data", ".txt")
	err := Import(ImportOptions{SnapshotPath: snap, SourcePath: src, Format: "xml"})
	if err == nil {
		t.Fatal("expected error")
	}
}
