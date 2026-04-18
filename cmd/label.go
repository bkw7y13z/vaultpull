package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/snapshot"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage human-readable labels for snapshot entries",
}

var labelAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Attach a label to a snapshot entry by checksum",
	RunE:  runLabelAdd,
}

var labelGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve the label for a snapshot entry",
	RunE:  runLabelGet,
}

var labelChecksum string
var labelText string
var labelSnapshotFile string

func init() {
	labelAddCmd.Flags().StringVar(&labelSnapshotFile, "snapshot", "snapshot.json", "Path to snapshot file")
	labelAddCmd.Flags().StringVar(&labelChecksum, "checksum", "", "Checksum of the entry to label")
	labelAddCmd.Flags().StringVar(&labelText, "label", "", "Label text to attach")

	labelGetCmd.Flags().StringVar(&labelSnapshotFile, "snapshot", "snapshot.json", "Path to snapshot file")
	labelGetCmd.Flags().StringVar(&labelChecksum, "checksum", "", "Checksum of the entry")

	labelCmd.AddCommand(labelAddCmd)
	labelCmd.AddCommand(labelGetCmd)
	rootCmd.AddCommand(labelCmd)
}

func runLabelAdd(cmd *cobra.Command, args []string) error {
	if labelChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if labelText == "" {
		return fmt.Errorf("--label is required")
	}
	if err := snapshot.Label(labelSnapshotFile, labelChecksum, labelText); err != nil {
		return fmt.Errorf("label: %w", err)
	}
	fmt.Printf("Labeled %s as %q\n", labelChecksum, labelText)
	return nil
}

func runLabelGet(cmd *cobra.Command, args []string) error {
	if labelChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	lbl, err := snapshot.GetLabel(labelSnapshotFile, labelChecksum)
	if err != nil {
		return fmt.Errorf("get label: %w", err)
	}
	fmt.Printf("Label for %s: %s\n", labelChecksum, lbl)
	return nil
}
