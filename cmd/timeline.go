package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Display snapshot history as an ordered timeline",
	RunE:  runTimeline,
}

func init() {
	timelineCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	timelineCmd.Flags().String("since", "", "Show entries after this time (RFC3339)")
	timelineCmd.Flags().String("until", "", "Show entries before this time (RFC3339)")
	timelineCmd.Flags().Bool("tagged", false, "Only show tagged entries")
	rootCmd.AddCommand(timelineCmd)
}

func runTimeline(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	sinceStr, _ := cmd.Flags().GetString("since")
	untilStr, _ := cmd.Flags().GetString("until")
	tagged, _ := cmd.Flags().GetBool("tagged")

	opts := snapshot.TimelineOptions{Tagged: tagged}

	if sinceStr != "" {
		t, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			return fmt.Errorf("invalid --since: %w", err)
		}
		opts.Since = t
	}
	if untilStr != "" {
		t, err := time.Parse(time.RFC3339, untilStr)
		if err != nil {
			return fmt.Errorf("invalid --until: %w", err)
		}
		opts.Until = t
	}

	entries, err := snapshot.Timeline(snapshotPath, opts)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No snapshot entries found.")
		return nil
	}

	for _, e := range entries {
		tag := e.Tag
		if tag == "" {
			tag = "(untagged)"
		}
		fmt.Printf("%s  %-20s  keys=%-4d  %s\n",
			e.CreatedAt.Format(time.RFC3339), tag, e.KeyCount, e.Checksum[:8])
	}
	return nil
}
