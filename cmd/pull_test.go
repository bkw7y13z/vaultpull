package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunPull_MissingConfig(t *testing.T) {
	configPath = "/nonexistent/path/vaultpull.toml"
	overwrite = false
	append = false

	err := runPull(nil, nil)
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestRunPull_InvalidConfig(t *testing.T) {
	tmp := t.TempDir()
	cfgFile := filepath.Join(tmp, "vaultpull.toml")

	// Write a config missing required fields
	content := []byte(`[vault]\noutput_path = "/tmp/out.env"\n`)
	if err := os.WriteFile(cfgFile, content, 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	configPath = cfgFile
	err := runPull(nil, nil)
	if err == nil {
		t.Fatal("expected error for invalid config, got nil")
	}
}

func TestExecute_NoArgs(t *testing.T) {
	// Ensure Execute does not panic when called with default (missing) config
	// It will call os.Exit(1) on error, so we only test that rootCmd is wired.
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}
	if rootCmd.Use != "vaultpull" {
		t.Errorf("expected Use=vaultpull, got %q", rootCmd.Use)
	}
}

func TestFlags_Defaults(t *testing.T) {
	f := rootCmd.PersistentFlags()

	c, err := f.GetString("config")
	if err != nil {
		t.Fatalf("config flag: %v", err)
	}
	if c != "vaultpull.toml" {
		t.Errorf("expected default config=vaultpull.toml, got %q", c)
	}

	ow, err := f.GetBool("overwrite")
	if err != nil {
		t.Fatalf("overwrite flag: %v", err)
	}
	if ow {
		t.Error("expected overwrite default=false")
	}

	ap, err := f.GetBool("append")
	if err != nil {
		t.Fatalf("append flag: %v", err)
	}
	if ap {
		t.Error("expected append default=false")
	}
}
