package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTOML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vaultpull.toml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTOML: %v", err)
	}
	return p
}

const validTOML = `
[vault]
address     = "http://127.0.0.1:8200"
token       = "root"
mount_path  = "secret"
secret_path = "myapp/prod"

[output]
file_path = ".env"
overwrite = true

[audit]
enabled  = true
log_file = "audit.log"

[filter]
prefix       = "APP_"
exclude_keys = ["SECRET_KEY"]
`

func TestLoad_Success(t *testing.T) {
	p := writeTOML(t, validTOML)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("unexpected address: %s", cfg.Vault.Address)
	}
	if cfg.Filter.Prefix != "APP_" {
		t.Errorf("unexpected prefix: %s", cfg.Filter.Prefix)
	}
	if len(cfg.Filter.ExcludeKeys) != 1 || cfg.Filter.ExcludeKeys[0] != "SECRET_KEY" {
		t.Errorf("unexpected exclude_keys: %v", cfg.Filter.ExcludeKeys)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/vaultpull.toml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	toml := `
[vault]
secret_path = "myapp/prod"
[output]
file_path = ".env"
`
	p := writeTOML(t, toml)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected validation error for missing address")
	}
}

func TestLoad_MissingOutputPath(t *testing.T) {
	toml := `
[vault]
address     = "http://127.0.0.1:8200"
secret_path = "myapp/prod"
[output]
`
	p := writeTOML(t, toml)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected validation error for missing output.file_path")
	}
}

func TestLoad_InvalidTOML(t *testing.T) {
	p := writeTOML(t, "this is not valid toml :::")
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected parse error for invalid TOML")
	}
}
