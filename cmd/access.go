package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "Manage snapshot access logs",
}

var accessRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record an access event for a snapshot entry",
	RunE:  runAccessRecord,
}

var accessListCmd = &cobra.Command{
	Use:   "list",
	Short: "List access events for a snapshot entry",
	RunE:  runAccessList,
}

func init() {
	accessRecordCmd.Flags().String("snapshot", "", "Path to snapshot file")
	accessRecordCmd.Flags().String("checksum", "", "Entry checksum")
	accessRecordCmd.Flags().String("by", "", "Who is accessing")
	accessRecordCmd.Flags().String("action", "", "Action performed (e.g. read, write)")
	accessRecordCmd.Flags().String("reason", "", "Optional reason")

	accessListCmd.Flags().String("snapshot", "", "Path to snapshot file")
	accessListCmd.Flags().String("checksum", "", "Entry checksum")

	accessCmd.AddCommand(accessRecordCmd, accessListCmd)
	rootCmd.AddCommand(accessCmd)
}

func runAccessRecord(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	by, _ := cmd.Flags().GetString("by")
	action, _ := cmd.Flags().GetString("action")
	reason, _ := cmd.Flags().GetString("reason")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if by == "" {
		return fmt.Errorf("--by is required")
	}
	if action == "" {
		return fmt.Errorf("--action is required")
	}

	if err := snapshot.RecordAccess(snap, checksum, by, action, reason); err != nil {
		return fmt.Errorf("recording access: %w", err)
	}
	fmt.Println("Access event recorded.")
	return nil
}

func runAccessList(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}

	entries, err := snapshot.GetAccessLog(snap, checksum)
	if err != nil {
		return fmt.Errorf("retrieving access log: %w", err)
	}
	if len(entries) == 0 {
		fmt.Println("No access events found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ACCESSED_BY\tACTION\tAT\tREASON")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.AccessedBy, e.Action, e.At.Format("2006-01-02 15:04:05"), e.Reason)
	}
	return w.Flush()
}
