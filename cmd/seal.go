package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"vaultpull/internal/snapshot"
)

var sealCmd = &cobra.Command{
	Use:   "seal",
	Short: "Seal or inspect sealed snapshot entries",
}

var sealAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Seal a snapshot entry by checksum",
	RunE:  runSealAdd,
}

var sealInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get seal info for a checksum",
	RunE:  runSealInfo,
}

var (
	sealPath     string
	sealChecksum string
	sealedBy     string
	sealReason   string
)

func init() {
	sealAddCmd.Flags().StringVar(&sealPath, "seal-path", "seals.json", "Path to seal store file")
	sealAddCmd.Flags().StringVar(&sealChecksum, "checksum", "", "Checksum of the snapshot entry to seal")
	sealAddCmd.Flags().StringVar(&sealedBy, "by", "", "Identity of who is sealing the entry")
	sealAddCmd.Flags().StringVar(&sealReason, "reason", "", "Reason for sealing")

	sealInfoCmd.Flags().StringVar(&sealPath, "seal-path", "seals.json", "Path to seal store file")
	sealInfoCmd.Flags().StringVar(&sealChecksum, "checksum", "", "Checksum to look up")

	sealCmd.AddCommand(sealAddCmd, sealInfoCmd)
	rootCmd.AddCommand(sealCmd)
}

func runSealAdd(cmd *cobra.Command, args []string) error {
	if sealChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	if sealedBy == "" {
		return fmt.Errorf("--by is required")
	}
	if sealReason == "" {
		return fmt.Errorf("--reason is required")
	}
	if err := snapshot.Seal(sealPath, sealChecksum, sealedBy, sealReason); err != nil {
		return fmt.Errorf("seal failed: %w", err)
	}
	fmt.Printf("Sealed entry %s\n", sealChecksum)
	return nil
}

func runSealInfo(cmd *cobra.Command, args []string) error {
	if sealChecksum == "" {
		return fmt.Errorf("--checksum is required")
	}
	rec, err := snapshot.GetSeal(sealPath, sealChecksum)
	if err != nil {
		return fmt.Errorf("seal info: %w", err)
	}
	fmt.Printf("Checksum : %s\n", rec.Checksum)
	fmt.Printf("Sealed By: %s\n", rec.SealedBy)
	fmt.Printf("Sealed At: %s\n", rec.SealedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Reason   : %s\n", rec.Reason)
	return nil
}
