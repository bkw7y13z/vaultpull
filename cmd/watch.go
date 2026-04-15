package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/snapshot"
	"vaultpull/internal/vault"
)

var watchInterval time.Duration
var watchMaxCycles int

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Poll Vault for secret changes and report diffs",
	RunE:  runWatch,
}

func init() {
	watchCmd.Flags().StringP("config", "c", "vaultpull.toml", "Path to config file")
	watchCmd.Flags().DurationVar(&watchInterval, "interval", 30*time.Second, "Poll interval")
	watchCmd.Flags().IntVar(&watchMaxCycles, "max-cycles", 0, "Max poll cycles (0 = infinite)")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("watch: config error: %w", err)
	}

	client, err := vault.NewClient(cfg.Address, cfg.Token)
	if err != nil {
		return fmt.Errorf("watch: vault client error: %w", err)
	}

	fetcher := vault.NewFetcher(client)

	secretsFn := func() (map[string]string, error) {
		raw, fetchErr := fetcher.FetchSecrets(cfg.SecretPath)
		if fetchErr != nil {
			return nil, fetchErr
		}
		return vault.FilterSecrets(raw, cfg.Namespace, cfg.ExcludeKeys), nil
	}

	fmt.Fprintf(os.Stdout, "Watching %s every %s\n", cfg.SecretPath, watchInterval)

	return snapshot.Watch(secretsFn, snapshot.WatchOptions{
		SnapshotPath: cfg.SnapshotPath,
		Interval:     watchInterval,
		MaxCycles:    watchMaxCycles,
		OnChange: func(r snapshot.WatchResult) {
			fmt.Fprintf(os.Stdout, "[%s] CHANGED checksum=%s\n", r.Timestamp.Format(time.RFC3339), r.Checksum)
			if r.Diff != nil {
				for _, k := range r.Diff.Added {
					fmt.Fprintf(os.Stdout, "  + %s\n", k)
				}
				for _, k := range r.Diff.Removed {
					fmt.Fprintf(os.Stdout, "  - %s\n", k)
				}
				for _, k := range r.Diff.Changed {
					fmt.Fprintf(os.Stdout, "  ~ %s\n", k)
				}
			}
		},
		OnNoChange: func(r snapshot.WatchResult) {
			fmt.Fprintf(os.Stdout, "[%s] no change checksum=%s\n", r.Timestamp.Format(time.RFC3339), r.Checksum)
		},
	})
}
