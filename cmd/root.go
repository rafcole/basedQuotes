/*
Copyright © 2024 Richard Cole
*/
package cmd

import (
	snapshotpkg "cryptoSnapShot/snapshot"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cryptoSnapShot [venue] [crypto pair]",
	Short: "An application to generate and store OHLCV snapshots from various venues",
	Long:  `Usage: cryptoSnapShot [venue] [quote]/[base]`,
	Args:  cobra.MinimumNArgs(2),
	Run:   snapshotpkg.TakeSnapShot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
