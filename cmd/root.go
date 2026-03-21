package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/config"
	"github.com/nopecho/claude-television/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ctv",
	Short: "Claude Code TUI dashboard",
	Long:  "claude-television — A TUI dashboard for exploring your Claude Code configuration at a glance.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		home, _ := os.UserHomeDir()
		claudeHome := filepath.Join(home, ".claude")
		configDir := config.ConfigDir()
		regPath := filepath.Join(configDir, "channels.json")

		reg, err := channel.LoadRegistry(regPath)
		if err != nil {
			return fmt.Errorf("load registry: %w", err)
		}

		if cfg.Channels.AutoSync {
			discovered, err := channel.DiscoverChannels(filepath.Join(claudeHome, "projects"))
			if err != nil {
				return fmt.Errorf("discover channels: %w", err)
			}
			reg = channel.DiffSync(reg, discovered)
		}

		if len(reg.Channels) == 0 {
			fmt.Println("No channels found. Run 'ctv init' or use Claude Code in a project first.")
			return nil
		}

		applyConfig(reg, cfg)

		cacheTTL := config.ParseDuration(cfg.Channels.CacheTTL)
		cache := channel.NewCache(filepath.Join(configDir, "cache"), cacheTTL)

		var wg sync.WaitGroup
		for i := range reg.Channels {
			wg.Add(1)
			go func(ch *channel.Channel) {
				defer wg.Done()
				expected := channel.ExpectedFiles(ch)
				if cache.IsValid(ch.ID, expected) {
					entry, err := cache.Load(ch.ID)
					if err == nil && entry != nil {
						ch.Data = &entry.Data
						return
					}
				}
				data, mtimes, err := channel.LoadChannelData(ch, claudeHome)
				if err != nil {
					ch.Status = channel.StatusError
					return
				}
				ch.Data = data
				ch.Status = determineStatus(ch, data)
				cache.Save(&channel.CacheEntry{
					ChannelID:  ch.ID,
					Data:       *data,
					FileMtimes: mtimes,
				})
			}(&reg.Channels[i])
		}
		wg.Wait()

		channel.SaveRegistry(reg, regPath)

		return tui.RunChannels(reg.Channels, cfg)
	},
}

func applyConfig(reg *channel.Registry, cfg *config.Config) {
	pinSet := map[string]bool{}
	for _, p := range cfg.Channels.Pins {
		pinSet[p] = true
	}
	groupMap := map[string]string{}
	for group, ids := range cfg.Channels.Groups {
		for _, id := range ids {
			groupMap[id] = group
		}
	}
	for i := range reg.Channels {
		ch := &reg.Channels[i]
		if pinSet[ch.ID] || pinSet[ch.Name] {
			ch.Pinned = true
		}
		if g, ok := groupMap[ch.ID]; ok {
			ch.Group = g
		}
		if g, ok := groupMap[ch.Name]; ok {
			ch.Group = g
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
