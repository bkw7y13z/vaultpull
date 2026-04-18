package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a snapshot tag",
	Long:  `Rename an existing tag label in the snapshot file to a new name.`,
	RunE:  runRename,
}

func init() {
	renameCmd.Flags().String("snapshot", ".vaultpull.snapshot.json", "Path to snapshot file")
	renameCmd.Flags().String("old", "", "Existing tag label to rename")
	renameCmd.Flags().String("new", "", "New tag label")
	_ = renameCmd.MarkFlagRequired("old")
	_ = renameCmd.MarkFlagRequired("new")
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	oldTag, _ := cmd.Flags().GetString("old")
	newTag, _ := cmd.Flags().GetString("new")

	if snapshotPath == "" {
		return fmt.Errorf("snapshot path must not be empty")
	}

	if err := snapshot.RenameTag(snapshotPath, oldTag, newTag); err != nil {
		return fmt.Errorf("rename tag: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Tag %q renamed to %q\n", oldTag, newTag)
	return nil
}
