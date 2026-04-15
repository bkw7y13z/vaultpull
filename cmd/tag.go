package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var (
	tagSnapshotPath string
	tagChecksum     string
	tagLabel        string
	tagFind         bool
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag or look up a snapshot entry by checksum",
	RunE:  runTag,
}

func init() {
	tagCmd.Flags().StringVar(&tagSnapshotPath, "snapshot", ".vaultpull-snapshot.json", "Path to snapshot file")
	tagCmd.Flags().StringVar(&tagChecksum, "checksum", "", "Checksum of the snapshot entry to tag")
	tagCmd.Flags().StringVar(&tagLabel, "label", "", "Human-readable label to assign or search")
	tagCmd.Flags().BoolVar(&tagFind, "find", false, "Find entry by label instead of tagging")
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	if tagLabel == "" {
		return fmt.Errorf("--label is required")
	}

	if tagFind {
		entry, err := snapshot.FindByTag(tagSnapshotPath, tagLabel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}
		fmt.Printf("Tag:       %s\n", entry.Tag)
		fmt.Printf("Checksum:  %s\n", entry.Checksum)
		fmt.Printf("Timestamp: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05 UTC"))
		fmt.Printf("Keys:      %v\n", entry.Keys)
		return nil
	}

	if tagChecksum == "" {
		return fmt.Errorf("--checksum is required when tagging")
	}

	if err := snapshot.Tag(tagSnapshotPath, tagChecksum, tagLabel); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Printf("Tagged snapshot %s as %q\n", tagChecksum, tagLabel)
	return nil
}
