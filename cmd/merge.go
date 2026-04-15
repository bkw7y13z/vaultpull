package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/audit"
	"github.com/vaultpull/vaultpull/internal/config"
	"github.com/vaultpull/vaultpull/internal/env"
	"github.com/vaultpull/vaultpull/internal/vault"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge secrets from Vault into an existing .env file",
	Long:  "Fetch secrets from Vault and merge them into an existing .env file, preserving keys not present in Vault.",
	RunE:  runMerge,
}

func init() {
	mergeCmd.Flags().StringP("config", "c", "vaultpull.toml", "path to config file")
	mergeCmd.Flags().Bool("overwrite", false, "overwrite existing keys with Vault values")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	raw, err := vault.FetchSecrets(client, cfg.Vault.SecretPath)
	if err != nil {
		return fmt.Errorf("failed to fetch secrets: %w", err)
	}

	filtered := vault.FilterSecrets(raw, cfg.Filter.Prefix, cfg.Filter.ExcludeKeys)

	logger, _ := audit.NewLogger(cfg.Audit.LogPath)

	result, err := env.MergeEnvFile(cfg.Output.Path, filtered, overwrite)
	if err != nil {
		_ = logger.Log("merge", cfg.Output.Path, false)
		return fmt.Errorf("merge failed: %w", err)
	}

	_ = logger.Log("merge", cfg.Output.Path, true)

	sort.Strings(result.Added)
	sort.Strings(result.Updated)
	sort.Strings(result.Skipped)

	fmt.Printf("Merge complete: %d added, %d updated, %d skipped\n",
		len(result.Added), len(result.Updated), len(result.Skipped))

	if len(result.Added) > 0 {
		fmt.Printf("  Added:   %v\n", result.Added)
	}
	if len(result.Updated) > 0 {
		fmt.Printf("  Updated: %v\n", result.Updated)
	}
	if len(result.Skipped) > 0 {
		fmt.Printf("  Skipped: %v\n", result.Skipped)
	}

	return nil
}
