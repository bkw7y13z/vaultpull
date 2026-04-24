package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vaultpull/vaultpull/internal/snapshot"
)

var (
	deprecateSnapshotPath string
	deprecateChecksum     string
	deprecateReason       string
	deprecateBy           string
	deprecateSuggest      string
	deprecateGetFlag      bool
)

func init() {
	deprecateCmd := &cobra.Command{
		Use:   "deprecate",
		Short: "Mark or query a snapshot entry as deprecated",
		RunE:  runDeprecate,
	}
	deprecateCmd.Flags().StringVar(&deprecateSnapshotPath, "snapshot", "", "Path to snapshot file")
	deprecateCmd.Flags().StringVar(&deprecateChecksum, "checksum", "", "Checksum of the entry")
	deprecateCmd.Flags().StringVar(&deprecateReason, "reason", "", "Reason for deprecation")
	deprecateCmd.Flags().StringVar(&deprecateBy, "by", "", "Who is deprecating the entry")
	deprecateCmd.Flags().StringVar(&deprecateSuggest, "suggest", "", "Suggested replacement (optional)")
	deprecateCmd.Flags().BoolVar(&deprecateGetFlag, "get", false, "Retrieve deprecation info instead of setting")
	rootCmd.AddCommand(deprecateCmd)
}

func runDeprecate(cmd *cobra.Command, args []string) error {
	if deprecateSnapshotPath == "" {
		return fmt.Errorf("--snapshot is required")
	}
	if deprecateChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}

	if deprecateGetFlag {
		rec, err := snapshot.GetDeprecation(deprecateSnapshotPath, deprecateChecksum)
		if err != nil {
			return err
		}
		if rec == nil {
			fmt.Println("No deprecation record found.")
			return nil
		}
		fmt.Printf("Checksum:     %s\n", rec.Checksum)
		fmt.Printf("Reason:       %s\n", rec.Reason)
		fmt.Printf("Deprecated By: %s\n", rec.DeprecatedBy)
		fmt.Printf("Deprecated At: %s\n", rec.DeprecatedAt.Format("2006-01-02 15:04:05"))
		if rec.Suggest != "" {
			fmt.Printf("Suggestion:   %s\n", rec.Suggest)
		}
		return nil
	}

	if deprecateReason == "" {
		return fmt.Errorf("--reason is required")
	}
	if deprecateBy == "" {
		return fmt.Errorf("--by is required")
	}

	if err := snapshot.Deprecate(deprecateSnapshotPath, deprecateChecksum, deprecateReason, deprecateBy, deprecateSuggest); err != nil {
		return err
	}
	fmt.Printf("Entry %s marked as deprecated.\n", deprecateChecksum)
	return nil
}
