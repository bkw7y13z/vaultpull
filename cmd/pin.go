package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Pin or unpin a snapshot entry to protect it from pruning",
}

var pinAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Pin a snapshot entry by checksum",
	RunE:  runPinAdd,
}

var pinRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Unpin a snapshot entry by checksum",
	RunE:  runPinRemove,
}

var (
	pinSnapshotPath string
	pinChecksum     string
	pinReason       string
)

func init() {
	pinAddCmd.Flags().StringVar(&pinSnapshotPath, "snapshot", "snapshots.json", "Path to snapshot file")
	pinAddCmd.Flags().StringVar(&pinChecksum, "checksum", "", "Checksum of the entry to pin")
	pinAddCmd.Flags().StringVar(&pinReason, "reason", "", "Reason for pinning")

	pinRemoveCmd.Flags().StringVar(&pinSnapshotPath, "snapshot", "snapshots.json", "Path to snapshot file")
	pinRemoveCmd.Flags().StringVar(&pinChecksum, "checksum", "", "Checksum of the entry to unpin")

	pinCmd.AddCommand(pinAddCmd)
	pinCmd.AddCommand(pinRemoveCmd)
	rootCmd.AddCommand(pinCmd)
}

func runPinAdd(cmd *cobra.Command, args []string) error {
	if pinChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if pinReason == "" {
		return fmt.Errorf("--reason is required")
	}
	if err := snapshot.Pin(pinSnapshotPath, pinChecksum, pinReason); err != nil {
		return fmt.Errorf("pin: %w", err)
	}
	fmt.Printf("Pinned %s: %s\n", pinChecksum, pinReason)
	return nil
}

func runPinRemove(cmd *cobra.Command, args []string) error {
	if pinChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if err := snapshot.Unpin(pinSnapshotPath, pinChecksum); err != nil {
		return fmt.Errorf("unpin: %w", err)
	}
	fmt.Printf("Unpinned %s\n", pinChecksum)
	return nil
}
