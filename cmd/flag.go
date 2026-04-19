package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var flagCmd = &cobra.Command{
	Use:   "flag",
	Short: "Flag or list flagged keys in a snapshot",
}

var flagAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Flag a key in a snapshot entry",
	RunE:  runFlagAdd,
}

var flagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List flagged keys for a snapshot entry",
	RunE:  runFlagList,
}

func init() {
	flagAddCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	flagAddCmd.Flags().String("checksum", "", "Entry checksum")
	flagAddCmd.Flags().String("key", "", "Secret key to flag")
	flagAddCmd.Flags().String("value", "", "Value hint (optional)")
	flagAddCmd.Flags().String("by", "", "Who is flagging")
	flagAddCmd.Flags().String("reason", "", "Reason for flagging")

	flagListCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	flagListCmd.Flags().String("checksum", "", "Filter by checksum (optional)")

	flagCmd.AddCommand(flagAddCmd, flagListCmd)
	rootCmd.AddCommand(flagCmd)
}

func runFlagAdd(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	key, _ := cmd.Flags().GetString("key")
	value, _ := cmd.Flags().GetString("value")
	by, _ := cmd.Flags().GetString("by")
	reason, _ := cmd.Flags().GetString("reason")

	if err := snapshot.Flag(snap, checksum, key, value, by, reason); err != nil {
		return fmt.Errorf("flag: %w", err)
	}
	fmt.Printf("Flagged key %q in entry %s\n", key, checksum)
	return nil
}

func runFlagList(cmd *cobra.Command, _ []string) error {
	snap, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")

	flags, err := snapshot.GetFlags(snap, checksum)
	if err != nil {
		return fmt.Errorf("flag list: %w", err)
	}
	if len(flags) == 0 {
		fmt.Println("No flags found.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CHECKSUM\tKEY\tFLAGGED_BY\tREASON\tAT")
	for _, f := range flags {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			f.Checksum, f.Key, f.FlaggedBy, f.Reason, f.At.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}
