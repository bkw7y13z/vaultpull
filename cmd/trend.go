package cmd

import (
	"fmt"
	"time"

	"github.com/fvbommel/vaultpull/internal/snapshot"
	"github.com/spf13/cobra"
)

var trendCmd = &cobra.Command{
	Use:   "trend",
	Short: "Show key-count trend over snapshot history",
	RunE:  runTrend,
}

func init() {
	trendCmd.Flags().String("snapshot", "snapshots.json", "Path to snapshot file")
	trendCmd.Flags().Duration("since", 0, "Only include entries newer than this duration ago (e.g. 24h)")
	RootCmd.AddCommand(trendCmd)
}

func runTrend(cmd *cobra.Command, _ []string) error {
	path, _ := cmd.Flags().GetString("snapshot")
	sinceDur, _ := cmd.Flags().GetDuration("since")

	var since time.Time
	if sinceDur > 0 {
		since = time.Now().UTC().Add(-sinceDur)
	}

	points, err := snapshot.Trend(path, since)
	if err != nil {
		return fmt.Errorf("trend: %w", err)
	}

	if len(points) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No snapshot entries found.")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%-30s %-12s %-20s %s\n", "Timestamp", "Keys", "Checksum", "Tag")
	for _, p := range points {
		tag := p.Tag
		if tag == "" {
			tag = "-"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%-30s %-12d %-20s %s\n",
			p.At.Format(time.RFC3339), p.KeyCount, p.Checksum[:8], tag)
	}
	return nil
}
