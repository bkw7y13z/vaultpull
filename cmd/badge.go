package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vaultpull/vaultpull/internal/snapshot"
)

var badgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Manage status badges for snapshot entries",
}

var badgeSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Attach a status badge to a snapshot entry",
	RunE:  runBadgeSet,
}

var badgeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a badge for a snapshot entry",
	RunE:  runBadgeGet,
}

func init() {
	badgeSetCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	badgeSetCmd.Flags().String("checksum", "", "Entry checksum")
	badgeSetCmd.Flags().String("label", "", "Badge label (e.g. ci, deploy)")
	badgeSetCmd.Flags().String("status", "ok", "Badge status: ok, warning, error")
	badgeSetCmd.Flags().String("message", "", "Badge message")

	badgeGetCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	badgeGetCmd.Flags().String("checksum", "", "Entry checksum")

	badgeCmd.AddCommand(badgeSetCmd, badgeGetCmd)
	rootCmd.AddCommand(badgeCmd)
}

func runBadgeSet(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	label, _ := cmd.Flags().GetString("label")
	status, _ := cmd.Flags().GetString("status")
	message, _ := cmd.Flags().GetString("message")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if label == "" {
		return fmt.Errorf("--label is required")
	}
	if err := snapshot.SetBadge(snapshotPath, checksum, label, status, message); err != nil {
		return err
	}
	fmt.Printf("Badge set: [%s] %s — %s\n", status, label, message)
	return nil
}

func runBadgeGet(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	b, ok, err := snapshot.GetBadge(snapshotPath, checksum)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println("No badge found for checksum:", checksum)
		return nil
	}
	fmt.Printf("Label:   %s\nStatus:  %s\nMessage: %s\n", b.Label, b.Status, b.Message)
	return nil
}
