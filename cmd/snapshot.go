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
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Namespace,
			len(e.Keys),
			e.Checksum[:8],
		)
	}
	return w.Flush()
}
