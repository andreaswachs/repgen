package cmd

import (
	"fmt"

	"github.com/andreaswachs/repgen/internal"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints the application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("repgen version %s\n", internal.APP_VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
