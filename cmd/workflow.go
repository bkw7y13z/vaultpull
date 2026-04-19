package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"vaultpull/internal/snapshot"
)

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage workflows attached to snapshot entries",
}

var workflowCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a workflow for a snapshot entry",
	RunE:  runWorkflowCreate,
}

var workflowListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workflows for a snapshot entry",
	RunE:  runWorkflowList,
}

func init() {
	workflowCreateCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	workflowCreateCmd.Flags().String("checksum", "", "Entry checksum")
	workflowCreateCmd.Flags().String("by", "", "Creator name")
	workflowCreateCmd.Flags().StringSlice("steps", nil, "Comma-separated step names")

	workflowListCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	workflowListCmd.Flags().String("checksum", "", "Entry checksum")

	workflowCmd.AddCommand(workflowCreateCmd, workflowListCmd)
	rootCmd.AddCommand(workflowCmd)
}

func runWorkflowCreate(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	by, _ := cmd.Flags().GetString("by")
	steps, _ := cmd.Flags().GetStringSlice("steps")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if by == "" {
		return fmt.Errorf("--by is required")
	}
	if len(steps) == 0 {
		return fmt.Errorf("--steps is required")
	}

	if err := snapshot.CreateWorkflow(snapshotPath, checksum, by, steps); err != nil {
		return err
	}
	fmt.Printf("workflow created for %s with steps: %s\n", checksum, strings.Join(steps, ", "))
	return nil
}

func runWorkflowList(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}

	workflows, err := snapshot.GetWorkflows(snapshotPath, checksum)
	if err != nil {
		return err
	}
	if len(workflows) == 0 {
		fmt.Println("no workflows found")
		return nil
	}
	for _, w := range workflows {
		fmt.Printf("workflow by=%s created_at=%s steps=%d\n", w.CreatedBy, w.CreatedAt.Format("2006-01-02 15:04:05"), len(w.Steps))
		for _, s := range w.Steps {
			fmt.Printf("  - %s [%s]\n", s.Name, s.Status)
		}
	}
	return nil
}
