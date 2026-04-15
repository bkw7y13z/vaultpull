package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

func writeWatchConfig(t *testing.T, dir string, data map[string]interface{}) string {
	t.Helper()
	p := filepath.Join(dir, "vaultpull.toml")
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := toml.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunWatch_MissingConfig(t *testing.T) {
	cmd := watchCmd
	cmd.ResetFlags()
	init()

	err := cmd.Flags().Set("config", "/nonexistent/vaultpull.toml")
	if err != nil {
		t.Skip("flag not registered yet")
	}

	runErr := runWatch(watchCmd, nil)
	if runErr == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestRunWatch_InvalidConfig(t *testing.T) {
	dir := t.TempDir()
	p := writeWatchConfig(t, dir, map[string]interface{}{
		"token": "tok",
		// missing address and output_path
	})

	watchCmd.Flags().Set("config", p)
	err := runWatch(watchCmd, nil)
	if err == nil {
		t.Fatal("expected config validation error")
	}
}

func TestWatchCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "watch" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("watch command not registered on root")
	}
}

func TestWatchCmd_DefaultFlags(t *testing.T) {
	intervalFlag := watchCmd.Flags().Lookup("interval")
	if intervalFlag == nil {
		t.Fatal("interval flag not found")
	}
	if intervalFlag.DefValue != "30s" {
		t.Fatalf("expected default interval 30s, got %s", intervalFlag.DefValue)
	}

	maxFlag := watchCmd.Flags().Lookup("max-cycles")
	if maxFlag == nil {
		t.Fatal("max-cycles flag not found")
	}
	if maxFlag.DefValue != "0" {
		t.Fatalf("expected default max-cycles 0, got %s", maxFlag.DefValue)
	}
}
