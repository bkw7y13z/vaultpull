package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/username/vaultpull/internal/snapshot"
)

func init() {
	quotaSetCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	quotaSetCmd.Flags().String("checksum", "", "Entry checksum")
	quotaSetCmd.Flags().String("by", "", "Who is setting the quota")
	quotaSetCmd.Flags().Int("max-keys", 0, "Maximum number of keys allowed")
	quotaSetCmd.Flags().Int("max-size-kb", 0, "Maximum size in kilobytes")

	quotaGetCmd.Flags().String("snapshot", "snapshot.json", "Path to snapshot file")
	quotaGetCmd.Flags().String("checksum", "", "Entry checksum")

	quotaCmd.AddCommand(quotaSetCmd)
	quotaCmd.AddCommand(quotaGetCmd)
	rootCmd.AddCommand(quotaCmd)
}

var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Manage key and size quotas for snapshot entries",
}

var quotaSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a quota policy for a snapshot entry",
	RunE:  runQuotaSet,
}

var quotaGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the quota policy for a snapshot entry",
	RunE:  runQuotaGet,
}

func runQuotaSet(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	by, _ := cmd.Flags().GetString("by")
	maxKeys, _ := cmd.Flags().GetInt("max-keys")
	maxSizeKB, _ := cmd.Flags().GetInt("max-size-kb")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if by == "" {
		return fmt.Errorf("--by is required")
	}
	if err := snapshot.SetQuota(snapshotPath, checksum, by, maxKeys, maxSizeKB); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "quota set for %s\n", checksum)
	return nil
}

func runQuotaGet(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	p, ok, err := snapshot.GetQuota(snapshotPath, checksum)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Fprintf(os.Stdout, "no quota found for %s\n", checksum)
		return nil
	}
	fmt.Fprintf(os.Stdout, "checksum: %s\nmax_keys: %d\nmax_size_kb: %d\nset_by: %s\ncreated_at: %s\n",
		p.Checksum, p.MaxKeys, p.MaxSizeKB, p.SetBy, p.CreatedAt.Format("2006-01-02T15:04:05Z"))
	return nil
}
