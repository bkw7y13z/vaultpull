package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/vaultpull/internal/audit"
	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/env"
	"github.com/user/vaultpull/internal/vault"
)

func runPull(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	logger, err := audit.NewLogger(cfg.AuditLogPath)
	if err != nil {
		return fmt.Errorf("initializing audit logger: %w", err)
	}

	client, err := vault.NewClient(cfg.VaultAddress, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	raw, err := vault.FetchSecrets(client, cfg.SecretPath)
	if err != nil {
		_ = logger.Log(cfg.SecretPath, false, err.Error())
		return fmt.Errorf("fetching secrets from %q: %w", cfg.SecretPath, err)
	}

	filtered := vault.FilterSecrets(raw, cfg.Prefix, cfg.ExcludeKeys)

	writeOpts := env.WriteOptions{
		Overwrite: overwrite,
		Append:    append,
	}
	if err := env.WriteEnvFile(cfg.OutputPath, filtered, writeOpts); err != nil {
		_ = logger.Log(cfg.SecretPath, false, err.Error())
		return fmt.Errorf("writing env file to %q: %w", cfg.OutputPath, err)
	}

	if err := logger.Log(cfg.SecretPath, true, ""); err != nil {
		fmt.Printf("warning: audit log failed: %v\n", err)
	}

	fmt.Printf("✓ wrote %d secret(s) to %s\n", len(filtered), cfg.OutputPath)
	return nil
}
