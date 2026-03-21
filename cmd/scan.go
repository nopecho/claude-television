package cmd

import (
	"fmt"

	"github.com/nopecho/claude-television/internal/config"
	"github.com/spf13/cobra"
)

var (
	scanList   bool
	scanRemove string
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Manage project scan paths",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if scanList {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if len(cfg.Scan.Roots) == 0 {
				fmt.Println("No scan paths registered. Use: ctv scan <path>")
				return nil
			}
			for _, r := range cfg.Scan.Roots {
				fmt.Println(r)
			}
			return nil
		}

		if scanRemove != "" {
			if _, err := config.Load(); err != nil {
				return err
			}
			if err := config.RemoveScanRoot(scanRemove); err != nil {
				return err
			}
			fmt.Printf("Removed: %s\n", scanRemove)
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("path required. Usage: ctv scan <path>")
		}

		if _, err := config.Load(); err != nil {
			return err
		}
		if err := config.AddScanRoot(args[0]); err != nil {
			return err
		}
		fmt.Printf("Added: %s\n", args[0])
		return nil
	},
}

func init() {
	scanCmd.Flags().BoolVar(&scanList, "list", false, "List registered scan paths")
	scanCmd.Flags().StringVar(&scanRemove, "remove", "", "Remove a scan path")
	rootCmd.AddCommand(scanCmd)
}
