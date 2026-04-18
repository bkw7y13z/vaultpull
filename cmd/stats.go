package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show aggregate statistics for a snapshot file",
	RunE:  runStats,
}

func init() {
	statsCmd.Flags().String("snapshot", "vaultpull.snapshot.json", "Path to snapshot file")
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, _ []string) error {
	path, err := cmd.Flags().GetString("snapshot")
	if err != nil {
		return err
	}
	if path == "" {
		return fmt.Errorf("--snapshot is required")
	}

	s, err := snapshot.ComputeStats(path)
	if err != nil {
		return fmt.Errorf("stats: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Total entries : %d\n", s.TotalEntries)
	fmt.Fprintf(cmd.OutOrStdout(), "Unique keys   : %d\n", s.UniqueKeys)
	fmt.Fprintf(cmd.OutOrStdout(), "Tagged        : %d\n", s.TaggedEntries)
	fmt.Fprintf(cmd.OutOrStdout(), "Pinned        : %d\n", s.PinnedEntries)
	if s.TotalEntries > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Oldest        : %s\n", s.OldestAt.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(cmd.OutOrStdout(), "Newest        : %s\n", s.NewestAt.Format("2006-01-02 15:04:05"))
	}
	if len(s.TopKeys) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Top keys      : %v\n", s.TopKeys)
	}
	return nil
}
