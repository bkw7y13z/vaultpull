package snapshot_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"vaultpull/internal/snapshot"
)

func TestImport_Integration_EnvThenJSON(t *testing.T) {
	dir := t.TempDir()
	snap := filepath.Join(dir, "snap.json")

	// First import: env file
	envFile := filepath.Join(dir, "secrets.env")
	os.WriteFile(envFile, []byte("DB_PASS=secret\nAPI_KEY=key123\n"), 0600)

	if err := snapshot.Import(snapshot.ImportOptions{
		SnapshotPath: snap,
		SourcePath:   envFile,
		Format:       "env",
		Tag:          "v1",
	}); err != nil {
		t.Fatalf("first import: %v", err)
	}

	// Second import: json file with different secrets
	secrets := map[string]string{"NEW_KEY": "newval"}
	b, _ := json.Marshal(secrets)
	jsonFile := filepath.Join(dir, "secrets.json")
	os.WriteFile(jsonFile, b, 0600)

	if err := snapshot.Import(snapshot.ImportOptions{
		SnapshotPath: snap,
		SourcePath:   jsonFile,
		Format:       "json",
		Tag:          "v2",
	}); err != nil {
		t.Fatalf("second import: %v", err)
	}

	entries, err := snapshot.Load(snap)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Tag != "v1" || entries[1].Tag != "v2" {
		t.Errorf("unexpected tags: %s, %s", entries[0].Tag, entries[1].Tag)
	}

	// Verify diff between the two
	diffs := snapshot.Diff(entries[0], entries[1])
	if len(diffs) == 0 {
		t.Error("expected diffs between imported snapshots")
	}
}
