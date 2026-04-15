package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpull/internal/snapshot"
)

var (
	pruneMaxAge  string
	pruneKeepN   int
	pruneDryRun  bool
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove old snapshot entries based on age or count",
	Long: `Prune trims the local snapshot history.

Examples:
  vaultpull prune --keep-last 10
  vaultpull prune --max-age 720h
  vaultpull prune --max-age 168h --keep-last 5 --dry-run`,
	RunE: runPrune,
}

func init() {
	pruneCmd.Flags().StringVar(&pruneMaxAge, "max-age", "", "remove entries older than this duration (e.g. 720h, 30m)")
	pruneCmd.Flags().IntVar(&pruneKeepN, "keep-last", 0, "always retain at least this many recent entries")
	pruneCmd.Flags().BoolVar(&pruneDryRun, "dry-run", false, "show what would be removed without modifying the file")
	rootCmd.AddCommand(pruneCmd)
}

func runPrune(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	if cfgPath == "" {
		cfgPath = "vaultpull.toml"
	}

	cfg, err := loadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	opts := snapshot.PruneOptions{KeepLast: pruneKeepN}
	if pruneMaxAge != "" {
		d, err := time.ParseDuration(pruneMaxAge)
		if err != nil {
			return fmt.Errorf("invalid --max-age %q: %w", pruneMaxAge, err)
		}
		opts.MaxAge = d
	}

	if pruneKeepN == 0 && opts.MaxAge == 0 {
		return fmt.Errorf("at least one of --keep-last or --max-age must be specified")
	}

	if pruneDryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] would prune %s with options: keep-last=%d max-age=%s\n",
			cfg.SnapshotPath, pruneKeepN, pruneMaxAge)
		return nil
	}

	result, err := snapshot.Prune(cfg.SnapshotPath, opts)
	if err != nil {
		return fmt.Errorf("pruning snapshot: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Pruned snapshot: removed=%d retained=%d\n",
		result.Removed, result.Retained)
	return nil
}
