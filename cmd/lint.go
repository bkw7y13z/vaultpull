package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var lintSnapshotPath string

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Validate a snapshot file for common issues",
	Long: `Lint inspects a snapshot file and reports issues such as duplicate
checksums, missing timestamps, empty key lists, and blank key names.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLint(lintSnapshotPath)
	},
}

func init() {
	lintCmd.Flags().StringVar(&lintSnapshotPath, "snapshot", ".vaultpull.snapshot.json", "Path to the snapshot file")
	rootCmd.AddCommand(lintCmd)
}

func runLint(path string) error {
	result, err := snapshot.Lint(path)
	if err != nil {
		return fmt.Errorf("lint failed: %w", err)
	}

	if !result.HasIssues() {
		fmt.Println("✔ snapshot is valid — no issues found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "✖ found %d issue(s):\n", len(result.Issues))
	for _, issue := range result.Issues {
		fmt.Fprintf(os.Stderr, "  %s\n", issue)
	}
	return fmt.Errorf("snapshot lint failed with %d issue(s)", len(result.Issues))
}
