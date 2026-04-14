package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds the full vaultpull configuration.
type Config struct {
	Vault   VaultConfig   `toml:"vault"`
	Output  OutputConfig  `toml:"output"`
	Audit   AuditConfig   `toml:"audit"`
	Filter  FilterConfig  `toml:"filter"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string `toml:"address"`
	Token     string `toml:"token"`
	MountPath string `toml:"mount_path"`
	SecretPath string `toml:"secret_path"`
}

// OutputConfig controls .env file output behaviour.
type OutputConfig struct {
	FilePath  string `toml:"file_path"`
	Overwrite bool   `toml:"overwrite"`
	Append    bool   `toml:"append"`
}

// AuditConfig controls audit logging.
type AuditConfig struct {
	Enabled  bool   `toml:"enabled"`
	LogFile  string `toml:"log_file"`
}

// FilterConfig controls namespace / key filtering.
type FilterConfig struct {
	Prefix      string   `toml:"prefix"`
	ExcludeKeys []string `toml:"exclude_keys"`
}

// Load reads a TOML config file from the given path.
// If path is empty it falls back to "vaultpull.toml" in the
// current directory.
func Load(path string) (*Config, error) {
	if path == "" {
		path = "vaultpull.toml"
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("config file not found: " + path)
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate performs basic sanity checks on a loaded Config.
func validate(cfg *Config) error {
	if cfg.Vault.Address == "" {
		return errors.New("vault.address is required")
	}
	if cfg.Vault.SecretPath == "" {
		return errors.New("vault.secret_path is required")
	}
	if cfg.Output.FilePath == "" {
		return errors.New("output.file_path is required")
	}
	return nil
}
