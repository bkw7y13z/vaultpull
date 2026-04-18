package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"vaultpull/internal/snapshot"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Manage snapshot audit log",
}

var auditLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Display audit log for snapshots",
	RunE:  runAuditLog,
}

var auditRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record an audit event",
	RunE:  runAuditRecord,
}

func init() {
	auditLogCmd.Flags().String("snapshot", "snapshots.json", "Path to snapshot file")
	auditRecordCmd.Flags().String("snapshot", "snapshots.json", "Path to snapshot file")
	auditRecordCmd.Flags().String("action", "", "Action name")
	auditRecordCmd.Flags().String("checksum", "", "Snapshot checksum")
	auditRecordCmd.Flags().String("actor", "", "Actor performing the action")
	auditRecordCmd.Flags().String("detail", "", "Additional detail")
	auditCmd.AddCommand(auditLogCmd)
	auditCmd.AddCommand(auditRecordCmd)
	rootCmd.AddCommand(auditCmd)
}

func runAuditLog(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	events, err := snapshot.GetAuditLog(snapshotPath)
	if err != nil {
		return fmt.Errorf("audit log: %w", err)
	}
	if len(events) == 0 {
		fmt.Println("No audit events found.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTION\tCHECKSUM\tACTOR\tDETAIL")
	for _, e := range events {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Action, e.Checksum, e.Actor, e.Detail)
	}
	return w.Flush()
}

func runAuditRecord(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	action, _ := cmd.Flags().GetString("action")
	checksum, _ := cmd.Flags().GetString("checksum")
	actor, _ := cmd.Flags().GetString("actor")
	detail, _ := cmd.Flags().GetString("detail")
	if err := snapshot.RecordAudit(snapshotPath, action, checksum, actor, detail); err != nil {
		return fmt.Errorf("record audit: %w", err)
	}
	fmt.Println("Audit event recorded.")
	return nil
}
