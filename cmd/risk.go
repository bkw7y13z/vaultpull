package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/internal/snapshot"
)

func init() {
	riskCmd := &cobra.Command{
		Use:   "risk",
		Short: "Manage risk assessments for snapshot entries",
	}

	assessCmd := &cobra.Command{
		Use:   "assess",
		Short: "Assess the risk level of a snapshot entry",
		RunE:  runRiskAssess,
	}
	assessCmd.Flags().String("snapshot", "", "Path to snapshot file")
	assessCmd.Flags().String("checksum", "", "Entry checksum")
	assessCmd.Flags().String("level", "", "Risk level: low, medium, high, critical")
	assessCmd.Flags().String("reason", "", "Reason for the risk assessment")
	assessCmd.Flags().String("by", "", "Who is performing the assessment")

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get the risk assessment for a snapshot entry",
		RunE:  runRiskGet,
	}
	getCmd.Flags().String("snapshot", "", "Path to snapshot file")
	getCmd.Flags().String("checksum", "", "Entry checksum")

	riskCmd.AddCommand(assessCmd, getCmd)
	rootCmd.AddCommand(riskCmd)
}

func runRiskAssess(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	level, _ := cmd.Flags().GetString("level")
	reason, _ := cmd.Flags().GetString("reason")
	by, _ := cmd.Flags().GetString("by")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if level == "" {
		return fmt.Errorf("--level is required")
	}
	if reason == "" {
		return fmt.Errorf("--reason is required")
	}
	if by == "" {
		return fmt.Errorf("--by is required")
	}

	if err := snapshot.AssessRisk(snapshotPath, checksum, reason, by, snapshot.RiskLevel(level)); err != nil {
		return fmt.Errorf("assess risk: %w", err)
	}
	fmt.Printf("Risk assessed: %s → %s\n", checksum[:8], level)
	return nil
}

func runRiskGet(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}

	entry, found, err := snapshot.GetRisk(snapshotPath, checksum)
	if err != nil {
		return fmt.Errorf("get risk: %w", err)
	}
	if !found {
		fmt.Println("No risk assessment found.")
		return nil
	}
	fmt.Printf("Level:       %s\n", entry.Level)
	fmt.Printf("Reason:      %s\n", entry.Reason)
	fmt.Printf("Assessed by: %s\n", entry.AssessedBy)
	fmt.Printf("Assessed at: %s\n", entry.AssessedAt.Format("2006-01-02 15:04:05"))
	return nil
}
