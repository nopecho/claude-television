package cmd

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func loadAllChannels(channels []channel.Channel, claudeHome string, cache *channel.Cache) {
	var wg sync.WaitGroup
	for i := range channels {
		wg.Add(1)
		go func(ch *channel.Channel) {
			defer wg.Done()

			if cache != nil {
				expected := channel.ExpectedFiles(ch)
				if entry, valid := cache.LoadIfValid(ch.ID, expected); valid {
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

			if cache != nil {
				cache.Save(&channel.CacheEntry{
					ChannelID:  ch.ID,
					Data:       *data,
					FileMtimes: mtimes,
				})
			}
		}(&channels[i])
	}
	wg.Wait()
}

func loadGlobalChannel(claudeHome string) channel.Channel {
	ch := channel.Channel{
		ID:       "__global__",
		Path:     claudeHome,
		Name:     "Global Settings",
		Status:   channel.StatusHealthy,
		IsGlobal: true,
		Data:     &channel.ChannelData{},
	}

	if settings, err := claude.ParseSettings(filepath.Join(claudeHome, "settings.json")); err == nil {
		ch.Data.Settings = settings
		ch.Data.Hooks, _ = claude.ExtractHooks(settings, "global")
		ch.Data.MCPServers, _ = claude.ExtractMCPServers(settings, "global")
	}

	if claudeMD, err := claude.ParseClaudeMD(filepath.Join(claudeHome, "CLAUDE.md")); err == nil {
		ch.Data.ClaudeMD = claudeMD
	}

	installed, _ := claude.ParseInstalledPlugins(filepath.Join(claudeHome, "plugins", "installed_plugins.json"))
	var enabled map[string]bool
	if ch.Data.Settings != nil {
		enabled = ch.Data.Settings.EnabledPlugins
	}
	ch.Data.Plugins = claude.MergePluginData(installed, enabled)
	ch.Data.LocalSkills, _ = claude.ScanLocalSkills(filepath.Join(claudeHome, "skills"))

	ch.Data.HealthIssues = claude.CheckHealth(&claude.HealthInput{
		ClaudeMD:   ch.Data.ClaudeMD,
		Settings:   ch.Data.Settings,
		MCPServers: ch.Data.MCPServers,
		Hooks:      ch.Data.Hooks,
	})

	return ch
}

func determineStatus(ch *channel.Channel, data *channel.ChannelData) channel.ChannelStatus {
	if _, err := os.Stat(ch.Path); err != nil {
		return channel.StatusError
	}
	if data.Settings == nil && data.ClaudeMD == nil {
		return channel.StatusWarning
	}
	return channel.StatusHealthy
}
