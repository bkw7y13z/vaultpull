package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

func init() {
	evictCmd := &cobra.Command{
		Use:   "evict",
		Short: "Evict a snapshot entry or list eviction records",
	}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Evict a snapshot entry by checksum",
		RunE:  runEvictAdd,
	}
	addCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	addCmd.Flags().String("checksum", "", "Checksum of the entry to evict")
	addCmd.Flags().String("by", "", "Identity performing the eviction")
	addCmd.Flags().String("reason", "", "Reason for eviction")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all eviction records",
		RunE:  runEvictList,
	}
	listCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")

	evictCmd.AddCommand(addCmd, listCmd)
	rootCmd.AddCommand(evictCmd)
}

func runEvictAdd(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	by, _ := cmd.Flags().GetString("by")
	reason, _ := cmd.Flags().GetString("reason")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if by == "" {
		return fmt.Errorf("--by is required")
	}
	if reason == "" {
		return fmt.Errorf("--reason is required")
	}

	if err := snapshot.Evict(snapshotPath, checksum, by, reason); err != nil {
		return fmt.Errorf("evict: %w", err)
	}
	fmt.Printf("evicted %s\n", checksum)
	return nil
}

func runEvictList(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")

	records, err := snapshot.GetEvictions(snapshotPath)
	if err != nil {
		return fmt.Errorf("get evictions: %w", err)
	}
	if len(records) == 0 {
		fmt.Println("no eviction records found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CHECKSUM\tEVICTED_BY\tREASON\tEVICTED_AT")
	for _, r := range records {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			r.Checksum, r.EvictedBy, r.Reason, r.EvictedAt.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}
