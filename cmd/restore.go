package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/env"
	"vaultpull/internal/snapshot"
)

var (
	restoreSnapshotPath string
	restoreOutputPath   string
	restoreOverwrite    bool
)

func init() {
	restoreCmd := &cobra.Command{
		Use:   "restore <ref>",
		Short: "Restore secrets from a snapshot entry into a .env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runRestore,
	}

	restoreCmd.Flags().StringVar(&restoreSnapshotPath, "snapshot", "vaultpull.snapshot.json", "Path to snapshot file")
	restoreCmd.Flags().StringVar(&restoreOutputPath, "output", ".env", "Destination .env file path")
	restoreCmd.Flags().BoolVar(&restoreOverwrite, "overwrite", false, "Overwrite existing .env file")

	rootCmd.AddCommand(restoreCmd)
}

func runRestore(cmd *cobra.Command, args []string) error {
	ref := args[0]

	result, err := snapshot.Restore(restoreSnapshotPath, ref)
	if err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	if err := env.WriteEnvFile(restoreOutputPath, result.Secrets, restoreOverwrite); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Restored %d keys from snapshot %s (tag: %q) → %s\n",
		len(result.Keys), result.Checksum[:min(8, len(result.Checksum))], result.Tag, restoreOutputPath)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
