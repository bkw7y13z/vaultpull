package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import secrets from a file into a snapshot",
	RunE:  runImport,
}

func init() {
	importCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	importCmd.Flags().String("source", "", "Path to source file (.env or .json)")
	importCmd.Flags().String("format", "env", "Source format: env or json")
	importCmd.Flags().String("tag", "", "Optional tag for the imported entry")
	importCmd.Flags().Bool("overwrite", false, "Allow importing duplicate checksums")
	importCmd.MarkFlagRequired("source")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	source, _ := cmd.Flags().GetString("source")
	format, _ := cmd.Flags().GetString("format")
	tag, _ := cmd.Flags().GetString("tag")
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	if err := snapshot.Import(snapshot.ImportOptions{
		SnapshotPath: snapshotPath,
		SourcePath:   source,
		Format:       format,
		Tag:          tag,
		Overwrite:    overwrite,
	}); err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Imported secrets from %s into %s\n", source, snapshotPath)
	return nil
}
