package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var (
	annotateChecksum string
	annotateNote     string
	annotateGet      bool
)

func init() {
	annotateCmd := &cobra.Command{
		Use:   "annotate",
		Short: "Attach or retrieve a note on a snapshot entry",
		RunE:  runAnnotate,
	}
	annotateCmd.Flags().StringVar(&snapshotPath, "snapshot", "vaultpull.snapshot.json", "Path to snapshot file")
	annotateCmd.Flags().StringVar(&annotateChecksum, "checksum", "", "Checksum of the snapshot entry")
	annotateCmd.Flags().StringVar(&annotateNote, "note", "", "Note to attach")
	annotateCmd.Flags().BoolVar(&annotateGet, "get", false, "Retrieve existing annotation instead of setting one")
	rootCmd.AddCommand(annotateCmd)
}

func runAnnotate(cmd *cobra.Command, args []string) error {
	if annotateChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}

	if annotateGet {
		note, err := snapshot.GetAnnotation(snapshotPath, annotateChecksum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}
		if note == "" {
			fmt.Println("(no annotation)")
		} else {
			fmt.Println(note)
		}
		return nil
	}

	if annotateNote == "" {
		return fmt.Errorf("--note is required when setting an annotation")
	}

	if err := snapshot.Annotate(snapshotPath, annotateChecksum, annotateNote); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	fmt.Printf("Annotation saved for checksum %s\n", annotateChecksum)
	return nil
}
