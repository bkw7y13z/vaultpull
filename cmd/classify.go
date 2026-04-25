package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"vaultpull/internal/snapshot"
)

var (
	classifySnapshotPath string
	classifyCategory     string
	classifySetBy        string
	classifySensitive    bool
)

func init() {
	classifyCmd := &cobra.Command{
		Use:   "classify <checksum>",
		Short: "Classify a snapshot entry with a category and sensitivity flag",
		Args:  cobra.ExactArgs(1),
		RunE:  runClassify,
	}
	classifyCmd.Flags().StringVar(&classifySnapshotPath, "snapshot", "snapshot.json", "Path to snapshot file")
	classifyCmd.Flags().StringVar(&classifyCategory, "category", "", "Category label (e.g. confidential, public)")
	classifyCmd.Flags().StringVar(&classifySetBy, "set-by", "", "Actor performing the classification")
	classifyCmd.Flags().BoolVar(&classifySensitive, "sensitive", false, "Mark entry as containing sensitive data")
	_ = classifyCmd.MarkFlagRequired("category")
	_ = classifyCmd.MarkFlagRequired("set-by")

	getClassifyCmd := &cobra.Command{
		Use:   "get <checksum>",
		Short: "Get classification for a snapshot entry",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetClassify,
	}
	getClassifyCmd.Flags().StringVar(&classifySnapshotPath, "snapshot", "snapshot.json", "Path to snapshot file")

	parent := &cobra.Command{
		Use:   "classify",
		Short: "Manage snapshot entry classifications",
	}
	parent.AddCommand(classifyCmd)
	parent.AddCommand(getClassifyCmd)
	rootCmd.AddCommand(parent)
}

func runClassify(cmd *cobra.Command, args []string) error {
	checksum := args[0]
	if err := snapshot.Classify(classifySnapshotPath, checksum, classifyCategory, classifySetBy, classifySensitive); err != nil {
		return fmt.Errorf("classify: %w", err)
	}
	fmt.Printf("classified %s as %q (sensitive=%v) by %s\n", checksum, classifyCategory, classifySensitive, classifySetBy)
	return nil
}

func runGetClassify(cmd *cobra.Command, args []string) error {
	checksum := args[0]
	c, ok, err := snapshot.GetClassification(classifySnapshotPath, checksum)
	if err != nil {
		return fmt.Errorf("get classification: %w", err)
	}
	if !ok {
		fmt.Printf("no classification found for %s\n", checksum)
		return nil
	}
	fmt.Printf("checksum:  %s\ncategory:  %s\nsensitive: %v\nset_by:    %s\nset_at:    %s\n",
		c.Checksum, c.Category, c.Sensitive, c.SetBy, c.SetAt.Format("2006-01-02 15:04:05"))
	return nil
}
