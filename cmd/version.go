package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print ctv version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ctv %s (%s)\n", version, commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
