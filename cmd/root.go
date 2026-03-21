package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ctv",
	Short: "Claude Code TUI dashboard",
	Long:  "claude-television — A TUI dashboard for exploring your Claude Code configuration at a glance.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Task 5에서 TUI 실행 로직 추가
		fmt.Println("claude-television dashboard (coming soon)")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
