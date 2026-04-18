package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes attached to snapshot entries",
}

var noteAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a note to a snapshot entry by checksum",
	RunE:  runNoteAdd,
}

var noteGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the note for a snapshot entry by checksum",
	RunE:  runNoteGet,
}

func init() {
	noteAddCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	noteAddCmd.Flags().String("checksum", "", "Checksum of the entry")
	noteAddCmd.Flags().String("note", "", "Note text to attach")

	noteGetCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	noteGetCmd.Flags().String("checksum", "", "Checksum of the entry")

	noteCmd.AddCommand(noteAddCmd)
	noteCmd.AddCommand(noteGetCmd)
	rootCmd.AddCommand(noteCmd)
}

func runNoteAdd(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	note, _ := cmd.Flags().GetString("note")
	if checksum == "" {
		fmt.Fprintln(os.Stderr, "error: --checksum is required")
		os.Exit(1)
	}
	if note == "" {
		fmt.Fprintln(os.Stderr, "error: --note is required")
		os.Exit(1)
	}
	if err := snapshot.AddNote(snap, checksum, note); err != nil {
		return err
	}
	fmt.Printf("Note added to entry %s\n", checksum)
	return nil
}

func runNoteGet(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	if checksum == "" {
		fmt.Fprintln(os.Stderr, "error: --checksum is required")
		os.Exit(1)
	}
	note, err := snapshot.GetNote(snap, checksum)
	if err != nil {
		return err
	}
	fmt.Println(note)
	return nil
}
