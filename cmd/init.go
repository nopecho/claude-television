package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Discover and register channels from ~/.claude/projects/",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, _ := os.UserHomeDir()
		claudeHome := filepath.Join(home, ".claude")
		projectsDir := filepath.Join(claudeHome, "projects")
		configDir := config.ConfigDir()

		fmt.Println("Scanning ~/.claude/projects/ ...")
		discovered, err := channel.DiscoverChannels(projectsDir)
		if err != nil {
			return fmt.Errorf("discover channels: %w", err)
		}

		if len(discovered) == 0 {
			fmt.Println("No channels found. Make sure Claude Code has been used in at least one project.")
			return nil
		}

		reg := &channel.Registry{Channels: discovered}

		cacheTTL := config.ParseDuration("24h")
		cache := channel.NewCache(filepath.Join(configDir, "cache"), cacheTTL)

		loadAllChannels(reg.Channels, claudeHome, cache)

		regPath := filepath.Join(configDir, "channels.json")
		if err := channel.SaveRegistry(reg, regPath); err != nil {
			return fmt.Errorf("save registry: %w", err)
		}

		var healthy, warning, errCount int
		for _, ch := range reg.Channels {
			icon := statusChar(ch.Status)
			fmt.Printf("  %s %-20s %s\n", icon, ch.Name, ch.Path)
			switch ch.Status {
			case channel.StatusHealthy:
				healthy++
			case channel.StatusWarning:
				warning++
			case channel.StatusError:
				errCount++
			}
		}
		fmt.Printf("\n%d channels registered (healthy: %d, warning: %d, error: %d)\n",
			len(reg.Channels), healthy, warning, errCount)

		fmt.Println("\nTo enable directory navigation (Alt+Enter), add to your shell config:")
		fmt.Println(`  ctv() { local dir; dir="$(command ctv "$@")"; [ -d "$dir" ] && cd "$dir" || command ctv "$@"; }`)

		return nil
	},
}

func statusChar(s channel.ChannelStatus) string {
	switch s {
	case channel.StatusHealthy:
		return "●"
	case channel.StatusWarning:
		return "○"
	case channel.StatusError:
		return "✕"
	}
	return "?"
}

func init() {
	rootCmd.AddCommand(initCmd)
}
