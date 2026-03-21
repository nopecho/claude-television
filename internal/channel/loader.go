package channel

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nopecho/claude-television/internal/claude"
)

// LoadChannelData parses all data for a channel.
// claudeHome is ~/.claude (for global data like plugins/skills).
// Returns parsed data and file mtimes for caching.
func LoadChannelData(ch *Channel, claudeHome string) (*ChannelData, map[string]time.Time, error) {
	data := &ChannelData{}
	mtimes := map[string]time.Time{}

	claudeDir := filepath.Join(ch.Path, ".claude")

	// Settings
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if info, err := os.Stat(settingsPath); err == nil {
		mtimes[settingsPath] = info.ModTime()
		data.Settings, _ = claude.ParseSettings(settingsPath)
	}

	// Local Settings
	localSettingsPath := filepath.Join(claudeDir, "settings.local.json")
	if info, err := os.Stat(localSettingsPath); err == nil {
		mtimes[localSettingsPath] = info.ModTime()
		data.LocalSettings, _ = claude.ParseSettings(localSettingsPath)
	}

	// CLAUDE.md (root)
	claudeMDPath := filepath.Join(ch.Path, "CLAUDE.md")
	if info, err := os.Stat(claudeMDPath); err == nil {
		mtimes[claudeMDPath] = info.ModTime()
		data.ClaudeMD, _ = claude.ParseClaudeMD(claudeMDPath)
	}

	// Sub CLAUDE.md files
	data.SubClaudeMDs = scanSubClaudeMDs(ch.Path, claudeMDPath, mtimes)

	// Hooks (merge project + global if available)
	if data.Settings != nil {
		data.Hooks = claude.ExtractHooks(data.Settings, "project")
	}

	// MCP Servers from project settings
	if data.Settings != nil {
		data.MCPServers = claude.ExtractMCPServers(data.Settings, "project")
	}

	// Git info
	data.GitInfo = loadGitInfo(ch.Path)

	// Memory files (from ~/.claude/projects/{id}/memory/)
	if claudeHome != "" {
		memoryDir := filepath.Join(claudeHome, "projects", ch.ID, "memory")
		data.MemoryFiles, _ = claude.ScanMemoryFiles(memoryDir)
	}

	// Global data (plugins, skills) — loaded once, shared
	if claudeHome != "" {
		installed, _ := claude.ParseInstalledPlugins(filepath.Join(claudeHome, "plugins", "installed_plugins.json"))
		var enabled map[string]bool
		globalSettings, _ := claude.ParseSettings(filepath.Join(claudeHome, "settings.json"))
		if globalSettings != nil {
			enabled = globalSettings.EnabledPlugins
			// Merge global hooks
			globalHooks := claude.ExtractHooks(globalSettings, "global")
			data.Hooks = append(globalHooks, data.Hooks...)
			// Merge global MCP servers
			globalMCP := claude.ExtractMCPServers(globalSettings, "global")
			data.MCPServers = append(globalMCP, data.MCPServers...)
		}
		data.Plugins = claude.MergePluginData(installed, enabled)
		data.LocalSkills, _ = claude.ScanLocalSkills(filepath.Join(claudeHome, "skills"))
	}

	return data, mtimes, nil
}

// ExpectedFiles returns the list of files a channel should track for cache invalidation.
func ExpectedFiles(ch *Channel) []string {
	claudeDir := filepath.Join(ch.Path, ".claude")
	return []string{
		filepath.Join(claudeDir, "settings.json"),
		filepath.Join(claudeDir, "settings.local.json"),
		filepath.Join(ch.Path, "CLAUDE.md"),
	}
}

func scanSubClaudeMDs(projectDir, rootClaudeMD string, mtimes map[string]time.Time) []claude.ClaudeMD {
	var result []claude.ClaudeMD
	skipDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, ".worktrees": true,
		"dist": true, "build": true, ".next": true, "target": true,
		"__pycache__": true, ".venv": true, ".tox": true,
	}
	filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == "CLAUDE.md" && path != rootClaudeMD {
			if info, err := d.Info(); err == nil {
				mtimes[path] = info.ModTime()
			}
			if md, err := claude.ParseClaudeMD(path); err == nil {
				result = append(result, *md)
			}
		}
		return nil
	})
	return result
}

func loadGitInfo(projectPath string) *GitInfo {
	if _, err := os.Stat(filepath.Join(projectPath, ".git")); err != nil {
		return nil
	}
	info := &GitInfo{}
	if out, err := runGit(projectPath, "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
		info.Branch = strings.TrimSpace(out)
	}
	if out, err := runGit(projectPath, "log", "-1", "--format=%h|%s|%ci"); err == nil {
		parts := strings.SplitN(strings.TrimSpace(out), "|", 3)
		if len(parts) == 3 {
			info.LastCommit = parts[0]
			info.LastCommitMsg = parts[1]
			info.LastCommitAt = parts[2]
		}
	}
	if out, err := runGit(projectPath, "status", "--porcelain"); err == nil {
		trimmed := strings.TrimSpace(out)
		if trimmed == "" {
			info.DirtyFiles = 0
		} else {
			info.DirtyFiles = len(strings.Split(trimmed, "\n"))
		}
	}
	return info
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return string(out), err
}
