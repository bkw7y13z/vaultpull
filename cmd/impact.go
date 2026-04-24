package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var impactCmd = &cobra.Command{
	Use:   "impact",
	Short: "Show the blast radius of a secret key across snapshots",
	RunE:  runImpact,
}

func init() {
	impactCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	impactCmd.Flags().String("key", "", "Secret key to analyse (required)")
	_ = impactCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(impactCmd)
}

func runImpact(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	key, _ := cmd.Flags().GetString("key")

	report, err := snapshot.Impact(snapshotPath, key)
	if err != nil {
		return fmt.Errorf("impact analysis failed: %w", err)
	}

	fmt.Printf("Impact report for key: %s\n", report.Key)
	fmt.Printf("Total snapshot references: %d\n", report.TotalRefs)

	if len(report.Entries) == 0 {
		fmt.Println("No snapshots reference this key.")
		return nil
	}

	fmt.Printf("\n%-20s  %-30s  %s\n", "Checksum", "Last Seen", "Changed At")
	fmt.Println("----------------------------------------------------------------------")
	for _, e := range report.Entries {
		chk := e.Checksum
		if len(chk) > 16 {
			chk = chk[:16]
		}
		fmt.Printf("%-20s  %-30s  %s\n",
			chk,
			e.LastSeen.Format("2006-01-02 15:04:05"),
			e.ChangedAt.Format("2006-01-02 15:04:05"),
		)
	}
	return nil
}
