package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var mergeSnapshotCmd = &cobra.Command{
	Use:   "merge-snapshot",
	Short: "Merge entries from a source snapshot into a destination snapshot",
	RunE:  runMergeSnapshot,
}

func init() {
	mergeSnapshotCmd.Flags().String("dst", "", "Destination snapshot file (required)")
	mergeSnapshotCmd.Flags().String("src", "", "Source snapshot file (required)")
	_ = mergeSnapshotCmd.MarkFlagRequired("dst")
	_ = mergeSnapshotCmd.MarkFlagRequired("src")
	rootCmd.AddCommand(mergeSnapshotCmd)
}

func runMergeSnapshot(cmd *cobra.Command, _ []string) error {
	dst, _ := cmd.Flags().GetString("dst")
	src, _ := cmd.Flags().GetString("src")

	res, err := snapshot.Merge(dst, src)
	if err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Merge complete: added=%d skipped=%d total=%d\n",
		res.Added, res.Skipped, res.Total)
	return nil
}
