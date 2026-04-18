package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var expireCmd = &cobra.Command{
	Use:   "expire",
	Short: "Manage expiration times for snapshot entries",
}

var expireAddCmd = &cobra.Command{
	Use:   "set",
	Short: "Set an expiration time on a snapshot entry",
	RunE:  runExpireSet,
}

var expireGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the expiration time of a snapshot entry",
	RunE:  runExpireGet,
}

func init() {
	expireAddCmd.Flags().String("snapshot", "", "Path to snapshot file")
	expireAddCmd.Flags().String("checksum", "", "Checksum of the entry")
	expireAddCmd.Flags().String("set-by", "", "Who is setting the expiry")
	expireAddCmd.Flags().String("expires-at", "", "Expiry time in RFC3339 format")

	expireGetCmd.Flags().String("snapshot", "", "Path to snapshot file")
	expireGetCmd.Flags().String("checksum", "", "Checksum of the entry")

	expireCmd.AddCommand(expireAddCmd)
	expireCmd.AddCommand(expireGetCmd)
	rootCmd.AddCommand(expireCmd)
}

func runExpireSet(cmd *cobra.Command, _ []string) error {
	path, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	setBy, _ := cmd.Flags().GetString("set-by")
	expiresAtStr, _ := cmd.Flags().GetString("expires-at")

	if path == "" {
		return fmt.Errorf("--snapshot is required")
	}
	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if expiresAtStr == "" {
		return fmt.Errorf("--expires-at is required")
	}
	expiry, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		return fmt.Errorf("invalid --expires-at: %w", err)
	}
	if err := snapshot.SetExpiry(path, checksum, setBy, expiry); err != nil {
		return err
	}
	fmt.Printf("Expiry set for %s until %s\n", checksum, expiry.Format(time.RFC3339))
	return nil
}

func runExpireGet(cmd *cobra.Command, _ []string) error {
	path, _ := cmd.Flags().GetString("snapshot")
	checksum, _ := cmd.Flags().GetString("checksum")
	if path == "" {
		return fmt.Errorf("--snapshot is required")
	}
	if checksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	e, err := snapshot.GetExpiry(path, checksum)
	if err != nil {
		return err
	}
	if e == nil {
		fmt.Println("No expiry set for", checksum)
		return nil
	}
	fmt.Printf("Checksum: %s\nExpires:  %s\nSet by:   %s\n", e.Checksum, e.ExpiresAt.Format(time.RFC3339), e.SetBy)
	return nil
}
