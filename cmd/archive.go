package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive old snapshot entries to a directory",
	RunE:  runArchive,
}

func init() {
	archiveCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	archiveCmd.Flags().String("dir", "archive", "Directory to store archived entries")
	archiveCmd.Flags().Int("keep-last", 10, "Number of recent entries to keep in snapshot")
	rootCmd.AddCommand(archiveCmd)
}

func runArchive(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	archiveDir, _ := cmd.Flags().GetString("dir")
	keepLast, _ := cmd.Flags().GetInt("keep-last")

	if err := snapshot.Archive(snapshotPath, archiveDir, keepLast); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Printf("Archived old entries. Kept last %d in %s. Archive dir: %s\n", keepLast, snapshotPath, archiveDir)
	return nil
}
