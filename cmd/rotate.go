package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Archive the current snapshot file and optionally prune old archives",
	RunE:  runRotate,
}

func init() {
	rotateCmd.Flags().String("snapshot", ".vaultpull-snapshot.json", "Path to the snapshot file to rotate")
	rotateCmd.Flags().String("archive-dir", ".vaultpull-archives", "Directory to store rotated snapshot archives")
	rotateCmd.Flags().Int("max-age-days", 0, "Remove archived files older than N days (0 = disabled)")
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, _ []string) error {
	snapshotPath, err := cmd.Flags().GetString("snapshot")
	if err != nil {
		return err
	}
	archiveDir, err := cmd.Flags().GetString("archive-dir")
	if err != nil {
		return err
	}
	maxAgeDays, err := cmd.Flags().GetInt("max-age-days")
	if err != nil {
		return err
	}

	opts := snapshot.RotateOptions{
		ArchiveDir: archiveDir,
		MaxAgeDays: maxAgeDays,
	}

	if err := snapshot.Rotate(snapshotPath, opts); err != nil {
		fmt.Fprintf(os.Stderr, "rotate error: %v\n", err)
		return err
	}

	fmt.Printf("Snapshot rotated to archive directory: %s\n", archiveDir)
	return nil
}
