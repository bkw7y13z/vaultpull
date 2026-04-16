package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify integrity of a snapshot file",
	RunE:  runVerify,
}

func init() {
	verifyCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	rootCmd.AddCommand(verifyCmd)
}

func runVerify(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")

	results, err := snapshot.Verify(snapshotPath)
	if err != nil {
		return fmt.Errorf("verify failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("snapshot is empty — nothing to verify")
		return nil
	}

	hasErrors := false
	for _, r := range results {
		if r.Valid {
			fmt.Printf("[OK]   %s", r.Checksum)
		} else {
			fmt.Fprintf(os.Stderr, "[FAIL] %s — %s\n", r.Checksum, r.Reason)
			hasErrors = true
		}
		if r.Tag != "" {
			fmt.Printf(" (tag: %s)", r.Tag)
		}
		fmt.Println()
	}

	if hasErrors {
		return fmt.Errorf("snapshot verification failed")
	}

	fmt.Printf("\nverified %d entries — all OK\n", len(results))
	return nil
}
