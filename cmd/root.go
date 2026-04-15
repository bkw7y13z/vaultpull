package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configPath string
	overwrite  bool
	append     bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpull fetches secrets from a HashiCorp Vault instance
and writes them to a local .env file with optional namespace filtering
and audit logging.`,
	RunE: runPull,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&configPath, "config", "c", "vaultpull.toml",
		"path to the TOML config file",
	)
	rootCmd.PersistentFlags().BoolVar(
		&overwrite, "overwrite", false,
		"overwrite the output .env file if it already exists",
	)
	rootCmd.PersistentFlags().BoolVar(
		&append, "append", false,
		"append secrets to the output .env file instead of replacing it",
	)
}
