package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vaultpull/internal/snapshot"
)

var lifecycleCmd = &cobra.Command{
	Use:   "lifecycle",
	Short: "Manage snapshot entry lifecycle states",
}

var lifecycleSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set lifecycle state for a snapshot entry",
	RunE:  runLifecycleSet,
}

var lifecycleGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get lifecycle state for a snapshot entry",
	RunE:  runLifecycleGet,
}

func init() {
	lifecycleSetCmd.Flags().String("snapshot", "", "Path to snapshot file")
	lifecycleSetCmd.Flags().String("checksum", "", "Entry checksum")
	lifecycleSetCmd.Flags().String("state", "", "Lifecycle state (active|deprecated|retired)")
	lifecycleSetCmd.Flags().String("by", "", "Who is making the change")
	lifecycleSetCmd.Flags().String("reason", "", "Reason for state change")

	lifecycleGetCmd.Flags().String("snapshot", "", "Path to snapshot file")
	lifecycleGetCmd.Flags().String("checksum", "", "Entry checksum")

	lifecycleCmd.AddCommand(lifecycleSetCmd)
	lifecycleCmd.AddCommand(lifecycleGetCmd)
	rootCmd.AddCommand(lifecycleCmd)
}

func runLifecycleSet(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	state, _ := cmd.Flags().GetString("state")
	by, _ := cmd.Flags().GetString("by")
	reason, _ := cmd.Flags().GetString("reason")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if state == "" {
		return fmt.Errorf("--state is required")
	}
	if by == "" {
		return fmt.Errorf("--by is required")
	}
	if reason == "" {
		return fmt.Errorf("--reason is required")
	}

	err := snapshot.SetLifecycle(snapshotPath, checksum, snapshot.LifecycleState(state), by, reason)
	if err != nil {
		return err
	}
	fmt.Printf("lifecycle state set to %q for %s\n", state, checksum)
	return nil
}

func runLifecycleGet(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}

	e, ok, err := snapshot.GetLifecycle(snapshotPath, checksum)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println("no lifecycle state found")
		return nil
	}
	fmt.Printf("state: %s\nchanged_by: %s\nreason: %s\nchanged_at: %s\n", e.State, e.ChangedBy, e.Reason, e.ChangedAt.Format("2006-01-02 15:04:05"))
	return nil
}
