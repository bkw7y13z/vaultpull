package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var blameCmd = &cobra.Command{
	Use:   "blame",
	Short: "Show which snapshot last changed each secret key",
	RunE:  runBlame,
}

func init() {
	blameCmd.Flags().String("snapshot", "vaultpull.snapshot.json", "Path to snapshot file")
	rootCmd.AddCommand(blameCmd)
}

func runBlame(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	if snapshotPath == "" {
		return fmt.Errorf("--snapshot path is required")
	}

	entries, err := snapshot.Blame(snapshotPath)
	if err != nil {
		return fmt.Errorf("blame failed: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No keys found in snapshot.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tCHECKSUM\tTAG\tCHANGED AT")
	for _, e := range entries {
		tag := e.Tag
		if tag == "" {
			tag = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			e.Key,
			e.Checksum[:8],
			tag,
			e.ChangedAt.Format("2006-01-02 15:04:05"),
		)
	}
	return w.Flush()
}
