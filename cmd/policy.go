package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultpull/internal/snapshot"
)

func init() {
	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage retention and validation policies for snapshots",
	}

	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Set a policy for a snapshot entry",
		RunE:  runPolicySet,
	}
	setCmd.Flags().String("snapshot", "", "Path to snapshot file")
	setCmd.Flags().String("checksum", "", "Entry checksum")
	setCmd.Flags().String("created-by", "", "Who is setting the policy")
	setCmd.Flags().Int("max-age", 0, "Max age in days")
	setCmd.Flags().Int("min-keys", 0, "Minimum number of keys required")
	setCmd.Flags().Bool("require-tag", false, "Require a tag on the entry")

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get the policy for a snapshot entry",
		RunE:  runPolicyGet,
	}
	getCmd.Flags().String("snapshot", "", "Path to snapshot file")
	getCmd.Flags().String("checksum", "", "Entry checksum")

	policyCmd.AddCommand(setCmd, getCmd)
	rootCmd.AddCommand(policyCmd)
}

func runPolicySet(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	createdBy, _ := cmd.Flags().GetString("created-by")
	maxAge, _ := cmd.Flags().GetInt("max-age")
	minKeys, _ := cmd.Flags().GetInt("min-keys")
	requireTag, _ := cmd.Flags().GetBool("require-tag")

	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if createdBy == "" {
		return fmt.Errorf("--created-by is required")
	}

	p := snapshot.Policy{
		MaxAge:     maxAge,
		MinKeys:    minKeys,
		RequireTag: requireTag,
		CreatedBy:  createdBy,
	}
	if err := snapshot.SetPolicy(snapshotPath, checksum, p); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	fmt.Println("policy set")
	return nil
}

func runPolicyGet(cmd *cobra.Command, _ []string) error {
	snapshotPath, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	pol, found, err := snapshot.GetPolicy(snapshotPath, checksum)
	if err != nil {
		return err
	}
	if !found {
		fmt.Println("no policy found")
		return nil
	}
	fmt.Printf("max_age_days: %d\nmin_keys: %d\nrequire_tag: %v\ncreated_by: %s\n",
		pol.MaxAge, pol.MinKeys, pol.RequireTag, pol.CreatedBy)
	return nil
}
