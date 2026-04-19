package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"vaultpull/internal/snapshot"
)

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Manage alerts on snapshot entries",
}

var alertAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an alert to a snapshot entry",
	RunE:  runAlertAdd,
}

var alertListCmd = &cobra.Command{
	Use:   "list",
	Short: "List alerts for a snapshot entry",
	RunE:  runAlertList,
}

func init() {
	alertAddCmd.Flags().String("snapshot", "", "Path to snapshot file")
	alertAddCmd.Flags().String("checksum", "", "Entry checksum")
	alertAddCmd.Flags().String("message", "", "Alert message")
	alertAddCmd.Flags().String("by", "", "Creator name")
	alertAddCmd.Flags().String("severity", "info", "Severity: info, warning, critical")

	alertListCmd.Flags().String("snapshot", "", "Path to snapshot file")
	alertListCmd.Flags().String("checksum", "", "Entry checksum")

	alertCmd.AddCommand(alertAddCmd, alertListCmd)
	rootCmd.AddCommand(alertCmd)
}

func runAlertAdd(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	message, _ := cmd.Flags().GetString("message")
	by, _ := cmd.Flags().GetString("by")
	sev, _ := cmd.Flags().GetString("severity")
	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if message == "" {
		return fmt.Errorf("--message is required")
	}
	if by == "" {
		return fmt.Errorf("--by is required")
	}
	err := snapshot.AddAlert(snap, checksum, message, by, snapshot.AlertSeverity(sev))
	if err != nil {
		return err
	}
	fmt.Println("alert added")
	return nil
}

func runAlertList(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	alerts, err := snapshot.GetAlerts(snap, checksum)
	if err != nil {
		return err
	}
	if len(alerts) == 0 {
		fmt.Println("no alerts found")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SEVERITY\tMESSAGE\tCREATED_BY\tCREATED_AT")
	for _, a := range alerts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", a.Severity, a.Message, a.CreatedBy, a.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}
