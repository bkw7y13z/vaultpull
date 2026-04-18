package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"vaultpull/internal/snapshot"
)

var bookmarkCmd = &cobra.Command{
	Use:   "bookmark",
	Short: "Manage snapshot bookmarks",
}

var bookmarkAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a bookmark to a snapshot entry",
	RunE:  runBookmarkAdd,
}

var bookmarkGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a bookmark by label",
	RunE:  runBookmarkGet,
}

var bookmarkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all bookmarks",
	RunE:  runBookmarkList,
}

func init() {
	bookmarkAddCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	bookmarkAddCmd.Flags().String("label", "", "Bookmark label")
	bookmarkAddCmd.Flags().String("checksum", "", "Entry checksum to bookmark")
	bookmarkAddCmd.Flags().String("note", "", "Optional note")

	bookmarkGetCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	bookmarkGetCmd.Flags().String("label", "", "Bookmark label")

	bookmarkListCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")

	bookmarkCmd.AddCommand(bookmarkAddCmd, bookmarkGetCmd, bookmarkListCmd)
	rootCmd.AddCommand(bookmarkCmd)
}

func runBookmarkAdd(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	label, _ := cmd.Flags().GetString("label")
	checksum, _ := cmd.Flags().GetString("checksum")
	note, _ := cmd.Flags().GetString("note")
	if label == "" {
		return fmt.Errorf("--label is required")
	}
	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if err := snapshot.AddBookmark(snap, label, checksum, note); err != nil {
		return err
	}
	fmt.Printf("Bookmark %q added for checksum %s\n", label, checksum)
	return nil
}

func runBookmarkGet(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	label, _ := cmd.Flags().GetString("label")
	if label == "" {
		return fmt.Errorf("--label is required")
	}
	b, err := snapshot.GetBookmark(snap, label)
	if err != nil {
		return err
	}
	fmt.Printf("Label:    %s\nChecksum: %s\nNote:     %s\n", b.Label, b.Checksum, b.Note)
	return nil
}

func runBookmarkList(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	list, err := snapshot.ListBookmarks(snap)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Println("No bookmarks found.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "LABEL\tCHECKSUM\tNOTE")
	for _, b := range list {
		fmt.Fprintf(w, "%s\t%s\t%s\n", b.Label, b.Checksum, b.Note)
	}
	return w.Flush()
}
