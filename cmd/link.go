package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"vaultpull/internal/snapshot"
)

func init() {
	linkAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a directional link between two snapshot entries",
		RunE:  runLinkAdd,
	}
	linkAddCmd.Flags().String("snapshot", "", "Path to snapshot file")
	linkAddCmd.Flags().String("from", "", "Source checksum")
	linkAddCmd.Flags().String("to", "", "Target checksum")
	linkAddCmd.Flags().String("reason", "", "Reason for the link")
	linkAddCmd.Flags().String("by", "", "Who is creating the link")

	linkListCmd := &cobra.Command{
		Use:   "list",
		Short: "List links from a snapshot entry",
		RunE:  runLinkList,
	}
	linkListCmd.Flags().String("snapshot", "", "Path to snapshot file")
	linkListCmd.Flags().String("from", "", "Source checksum to list links for")

	linkCmd := &cobra.Command{
		Use:   "link",
		Short: "Manage directional links between snapshot entries",
	}
	linkCmd.AddCommand(linkAddCmd, linkListCmd)
	rootCmd.AddCommand(linkCmd)
}

func runLinkAdd(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	reason, _ := cmd.Flags().GetString("reason")
	by, _ := cmd.Flags().GetString("by")

	if from == "" {
		return fmt.Errorf("--from is required")
	}
	if to == "" {
		return fmt.Errorf("--to is required")
	}
	if err := snapshot.AddLink(snapshotPath, from, to, reason, by); err != nil {
		return err
	}
	fmt.Printf("linked %s -> %s (%s)\n", from, to, reason)
	return nil
}

func runLinkList(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	from, _ := cmd.Flags().GetString("from")
	if from == "" {
		return fmt.Errorf("--from is required")
	}
	links, err := snapshot.GetLinks(snapshotPath, from)
	if err != nil {
		return err
	}
	if len(links) == 0 {
		fmt.Println("no links found")
		return nil
	}
	for _, l := range links {
		fmt.Printf("%s -> %s | reason: %s | by: %s | at: %s\n",
			l.FromChecksum, l.ToChecksum, l.Reason, l.CreatedBy, l.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}
