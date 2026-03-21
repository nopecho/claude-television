package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/nopecho/claude-television/internal/claude"
	"github.com/nopecho/claude-television/internal/config"
	"github.com/nopecho/claude-television/internal/scanner"
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
		claudeDir := filepath.Join(home, ".claude")

		data := tui.DashboardData{}
		var wg sync.WaitGroup

		var installed map[string]claude.InstalledPlugin

		wg.Add(1)
		go func() {
			defer wg.Done()
			data.Settings, _ = claude.ParseSettings(filepath.Join(claudeDir, "settings.json"))
			data.LocalSettings, _ = claude.ParseSettings(filepath.Join(claudeDir, "settings.local.json"))
			data.ClaudeMD, _ = claude.ParseClaudeMD(filepath.Join(claudeDir, "CLAUDE.md"))
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			installed, _ = claude.ParseInstalledPlugins(filepath.Join(claudeDir, "plugins", "installed_plugins.json"))
			data.LocalSkills, _ = claude.ScanLocalSkills(filepath.Join(claudeDir, "skills"))
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			data.ProjectsMeta, _ = claude.ScanProjectsMeta(filepath.Join(claudeDir, "projects"))
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			data.Projects, _ = scanner.ScanProjects(cfg.Scan.Roots, cfg.Scan.Ignore)
		}()

		wg.Wait()

		var enabled map[string]bool
		if data.Settings != nil {
			enabled = data.Settings.EnabledPlugins
			data.Hooks = claude.ExtractHooks(data.Settings, "global")
		}
		data.Plugins = claude.MergePluginData(installed, enabled)

		return tui.Run(data)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
