package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/vaultpull/internal/snapshot"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock or unlock snapshot entries",
}

var lockAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Lock a snapshot entry by checksum",
	RunE:  runLockAdd,
}

var lockRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Unlock a snapshot entry by checksum",
	RunE:  runLockRemove,
}

var (
	lockSnapshotFile string
	lockChecksum     string
	lockLockedBy     string
	lockReason       string
)

func init() {
	lockAddCmd.Flags().StringVar(&lockSnapshotFile, "snapshot", "snapshot.json", "Path to snapshot file")
	lockAddCmd.Flags().StringVar(&lockChecksum, "checksum", "", "Checksum of the entry to lock")
	lockAddCmd.Flags().StringVar(&lockLockedBy, "by", "", "Identity of the locker")
	lockAddCmd.Flags().StringVar(&lockReason, "reason", "", "Reason for locking")

	lockRemoveCmd.Flags().StringVar(&lockSnapshotFile, "snapshot", "snapshot.json", "Path to snapshot file")
	lockRemoveCmd.Flags().StringVar(&lockChecksum, "checksum", "", "Checksum of the entry to unlock")

	lockCmd.AddCommand(lockAddCmd)
	lockCmd.AddCommand(lockRemoveCmd)
	rootCmd.AddCommand(lockCmd)
}

func runLockAdd(cmd *cobra.Command, args []string) error {
	if lockChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if lockLockedBy == "" {
		return fmt.Errorf("--by is required")
	}
	if lockReason == "" {
		return fmt.Errorf("--reason is required")
	}
	if err := snapshot.Lock(lockSnapshotFile, lockChecksum, lockLockedBy, lockReason); err != nil {
		return fmt.Errorf("lock failed: %w", err)
	}
	fmt.Printf("Entry %s locked by %s\n", lockChecksum, lockLockedBy)
	return nil
}

func runLockRemove(cmd *cobra.Command, args []string) error {
	if lockChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if err := snapshot.Unlock(lockSnapshotFile, lockChecksum); err != nil {
		return fmt.Errorf("unlock failed: %w", err)
	}
	fmt.Printf("Entry %s unlocked\n", lockChecksum)
	return nil
}
