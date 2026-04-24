package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nicholasgasior/vaultpull/internal/config"
	"github.com/nicholasgasior/vaultpull/internal/snapshot"
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Show the pull history recorded in the local snapshot file",
	RunE:  runSnapshot,
}

func init() {
	snapshotCmd.Flags().StringP("config", "c", "vaultpull.toml", "Path to config file")
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	snapPath := cfg.SnapshotPath
	if snapPath == "" {
		snapPath = ".vaultpull.snapshot.json"
	}

	s, err := snapshot.Load(snapPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	if len(s.Entries) == 0 {
		fmt.Println("No snapshot entries found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tNAMESPACE\tKEYS\tCHECKSUM")
	for _, e := range s.Entries {
		checksum := formatChecksum(e.Checksum)
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Namespace,
			len(e.Keys),
			checksum,
		)
	}
	return w.Flush()
}

// formatChecksum returns a shortened display version of a checksum string.
// If the checksum is shorter than 8 characters, the full value is returned
// to avoid a panic on slice bounds out of range.
func formatChecksum(checksum string) string {
	if len(checksum) < 8 {
		return checksum
	}
	return checksum[:8]
}
