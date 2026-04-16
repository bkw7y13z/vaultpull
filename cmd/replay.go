package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var replayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Walk snapshot history between two refs and show diffs",
	RunE:  runReplay,
}

func init() {
	replayCmd.Flags().String("snapshot", "vaultpull.snapshot.json", "Path to snapshot file")
	replayCmd.Flags().String("from", "", "Starting checksum or tag (required)")
	replayCmd.Flags().String("to", "", "Ending checksum or tag (default: latest)")
	replayCmd.Flags().Bool("dry", false, "Print steps without invoking handler")
	_ = replayCmd.MarkFlagRequired("from")
	rootCmd.AddCommand(replayCmd)
}

func runReplay(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	dry, _ := cmd.Flags().GetBool("dry")

	opts := snapshot.ReplayOptions{From: from, To: to, Dry: dry}

	step := 0
	err := snapshot.Replay(snapshotPath, opts, func(e snapshot.ReplayEvent) error {
		step++
		tag := e.Tag
		if tag == "" {
			tag = "(untagged)"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "[%d] %s  tag=%s  keys=%s\n",
			e.Index, e.Checksum[:8], tag, strings.Join(e.Keys, ","))
		if e.Diff != nil {
			if len(e.Diff.Added) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "    + added: %s\n", strings.Join(e.Diff.Added, ", "))
			}
			if len(e.Diff.Removed) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "    - removed: %s\n", strings.Join(e.Diff.Removed, ", "))
			}
			if len(e.Diff.Changed) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "    ~ changed: %s\n", strings.Join(e.Diff.Changed, ", "))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if dry {
		fmt.Fprintln(cmd.OutOrStdout(), "(dry run — no handler invoked)")
	}
	return nil
}
