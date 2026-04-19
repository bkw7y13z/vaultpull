package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var (
	promoteFromEnv    string
	promoteToEnv      string
	promotePromotedBy string
	promoteNote       string
	promoteChecksum   string
	promoteListMode   bool
)

func init() {
	promoteCmd := &cobra.Command{
		Use:   "promote",
		Short: "Record or list environment promotions for a snapshot entry",
		RunE:  runPromote,
	}
	promoteCmd.Flags().StringVar(&snapshotFile, "snapshot", "vaultpull.snap.json", "Path to snapshot file")
	promoteCmd.Flags().StringVar(&promoteChecksum, "checksum", "", "Entry checksum")
	promoteCmd.Flags().StringVar(&promoteFromEnv, "from", "", "Source environment")
	promoteCmd.Flags().StringVar(&promoteToEnv, "to", "", "Target environment")
	promoteCmd.Flags().StringVar(&promotePromotedBy, "by", "", "Who performed the promotion")
	promoteCmd.Flags().StringVar(&promoteNote, "note", "", "Optional note")
	promoteCmd.Flags().BoolVar(&promoteListMode, "list", false, "List promotions for checksum")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	if promoteChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if promoteListMode {
		records, err := snapshot.GetPromotions(snapshotFile, promoteChecksum)
		if err != nil {
			return fmt.Errorf("get promotions: %w", err)
		}
		if len(records) == 0 {
			fmt.Println("No promotions found.")
			return nil
		}
		for _, r := range records {
			fmt.Printf("[%s] %s → %s by %s", r.PromotedAt.Format("2006-01-02 15:04:05"), r.FromEnv, r.ToEnv, r.PromotedBy)
			if r.Note != "" {
				fmt.Printf(" (%s)", r.Note)
			}
			fmt.Println()
		}
		return nil
	}
	if promoteFromEnv == "" {
		return fmt.Errorf("--from is required")
	}
	if promoteToEnv == "" {
		return fmt.Errorf("--to is required")
	}
	if promotePromotedBy == "" {
		return fmt.Errorf("--by is required")
	}
	if err := snapshot.Promote(snapshotFile, promoteChecksum, promoteFromEnv, promoteToEnv, promotePromotedBy, promoteNote); err != nil {
		return fmt.Errorf("promote: %w", err)
	}
	fmt.Printf("Promoted %s from %s to %s\n", promoteChecksum, promoteFromEnv, promoteToEnv)
	return nil
}
