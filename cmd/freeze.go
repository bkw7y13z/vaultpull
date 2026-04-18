package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var freezeCmd = &cobra.Command{
	Use:   "freeze",
	Short: "Freeze or inspect frozen snapshot entries",
}

var freezeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Freeze a snapshot entry by checksum",
	RunE:  runFreezeAdd,
}

var freezeInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show freeze record for a checksum",
	RunE:  runFreezeInfo,
}

var (
	freezeSnapshotFile string
	freezeChecksum     string
	freezeFrozenBy     string
	freezeReason       string
)

func init() {
	freezeAddCmd.Flags().StringVar(&freezeSnapshotFile, "snapshot", "snapshot.json", "Path to snapshot file")
	freezeAddCmd.Flags().StringVar(&freezeChecksum, "checksum", "", "Checksum to freeze")
	freezeAddCmd.Flags().StringVar(&freezeFrozenBy, "by", "", "Who is freezing the entry")
	freezeAddCmd.Flags().StringVar(&freezeReason, "reason", "", "Reason for freezing")

	freezeInfoCmd.Flags().StringVar(&freezeSnapshotFile, "snapshot", "snapshot.json", "Path to snapshot file")
	freezeInfoCmd.Flags().StringVar(&freezeChecksum, "checksum", "", "Checksum to inspect")

	freezeCmd.AddCommand(freezeAddCmd)
	freezeCmd.AddCommand(freezeInfoCmd)
	rootCmd.AddCommand(freezeCmd)
}

func runFreezeAdd(cmd *cobra.Command, args []string) error {
	if freezeChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if freezeFrozenBy == "" {
		return fmt.Errorf("--by is required")
	}
	if freezeReason == "" {
		return fmt.Errorf("--reason is required")
	}
	if err := snapshot.Freeze(freezeSnapshotFile, freezeChecksum, freezeFrozenBy, freezeReason); err != nil {
		return err
	}
	fmt.Printf("Entry %s frozen by %s\n", freezeChecksum, freezeFrozenBy)
	return nil
}

func runFreezeInfo(cmd *cobra.Command, args []string) error {
	if freezeChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	r, ok, err := snapshot.GetFreeze(freezeSnapshotFile, freezeChecksum)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Printf("Entry %s is not frozen\n", freezeChecksum)
		return nil
	}
	fmt.Printf("Frozen by: %s\nReason:    %s\nFrozen at: %s\n", r.FrozenBy, r.Reason, r.FrozenAt.Format("2006-01-02 15:04:05"))
	return nil
}
