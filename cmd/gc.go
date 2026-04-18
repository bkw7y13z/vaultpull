package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "Garbage collect old snapshot entries",
	RunE:  runGC,
}

func init() {
	gcCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	gcCmd.Flags().Duration("max-age", 30*24*time.Hour, "Maximum age of entries to keep")
	gcCmd.Flags().Bool("keep-pinned", true, "Keep pinned entries regardless of age")
	gcCmd.Flags().Bool("keep-tagged", true, "Keep tagged entries regardless of age")
	gcCmd.Flags().Bool("dry-run", false, "Report what would be removed without modifying")
	rootCmd.AddCommand(gcCmd)
}

func runGC(cmd *cobra.Command, _ []string) error {
	path, _ := cmd.Flags().GetString("snapshot")
	maxAge, _ := cmd.Flags().GetDuration("max-age")
	keepPinned, _ := cmd.Flags().GetBool("keep-pinned")
	keepTagged, _ := cmd.Flags().GetBool("keep-tagged")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	result, err := snapshot.GC(snapshot.GCOptions{
		SnapshotPath: path,
		MaxAge:       maxAge,
		KeepPinned:   keepPinned,
		KeepTagged:   keepTagged,
		DryRun:       dryRun,
	})
	if err != nil {
		return fmt.Errorf("gc failed: %w", err)
	}

	if dryRun {
		fmt.Printf("[dry-run] would remove %d entries\n", len(result.Removed))
	} else {
		fmt.Printf("removed %d entries, kept %d\n", len(result.Removed), len(result.Kept))
	}
	for _, c := range result.Removed {
		fmt.Printf("  - %s\n", c)
	}
	return nil
}
