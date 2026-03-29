package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/config"
	"github.com/spf13/cobra"
)

type exportChannel struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Path         string               `json:"path"`
	Status       channel.ChannelStatus `json:"status"`
	IsGlobal     bool                 `json:"is_global"`
	SettingsSummary *settingsSummary  `json:"settings_summary,omitempty"`
	HealthIssues interface{}          `json:"health_issues,omitempty"`
	GitInfo      interface{}          `json:"git_info,omitempty"`
}

type settingsSummary struct {
	HookCount      int `json:"hook_count"`
	MCPServerCount int `json:"mcp_server_count"`
	PluginCount    int `json:"plugin_count"`
	SkillCount     int `json:"skill_count"`
}

var exportChannelName string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export channel data as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, _ := os.UserHomeDir()
		claudeHome := filepath.Join(home, ".claude")
		configDir := config.ConfigDir()
		regPath := filepath.Join(configDir, "channels.json")

		reg, err := channel.LoadRegistry(regPath)
		if err != nil {
			return fmt.Errorf("load registry: %w", err)
		}

		cacheTTL := config.ParseDuration("24h")
		cache := channel.NewCache(filepath.Join(configDir, "cache"), cacheTTL)
		loadAllChannels(reg.Channels, claudeHome, cache)

		globalCh := loadGlobalChannel(claudeHome)
		allChannels := append([]channel.Channel{globalCh}, reg.Channels...)

		if exportChannelName != "" {
			for _, ch := range allChannels {
				if strings.EqualFold(ch.Name, exportChannelName) || ch.ID == exportChannelName {
					return printJSON(toExportChannel(ch))
				}
			}
			return fmt.Errorf("channel not found: %s", exportChannelName)
		}

		out := make([]exportChannel, 0, len(allChannels))
		for _, ch := range allChannels {
			out = append(out, toExportChannel(ch))
		}
		return printJSON(out)
	},
}

func toExportChannel(ch channel.Channel) exportChannel {
	e := exportChannel{
		ID:       ch.ID,
		Name:     ch.Name,
		Path:     ch.Path,
		Status:   ch.Status,
		IsGlobal: ch.IsGlobal,
	}
	if ch.Data != nil {
		e.SettingsSummary = &settingsSummary{
			HookCount:      len(ch.Data.Hooks),
			MCPServerCount: len(ch.Data.MCPServers),
			PluginCount:    len(ch.Data.Plugins),
			SkillCount:     len(ch.Data.LocalSkills),
		}
		if len(ch.Data.HealthIssues) > 0 {
			e.HealthIssues = ch.Data.HealthIssues
		}
		if ch.Data.GitInfo != nil {
			e.GitInfo = ch.Data.GitInfo
		}
	}
	return e
}

func printJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func init() {
	exportCmd.Flags().StringVarP(&exportChannelName, "channel", "c", "", "Export a specific channel by name or ID")
	rootCmd.AddCommand(exportCmd)
}
