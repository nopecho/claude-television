package channel

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nopecho/claude-television/internal/claude"
	"golang.org/x/sync/errgroup"
)

type channelLoader struct {
	ch         *Channel
	claudeHome string
}

func (l *channelLoader) load() (*ChannelData, map[string]time.Time, error) {
	data := &ChannelData{}
	mtimes := map[string]time.Time{}

	claudeDir := filepath.Join(l.ch.Path, ".claude")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	localSettingsPath := filepath.Join(claudeDir, "settings.local.json")
	claudeMDPath := filepath.Join(l.ch.Path, "CLAUDE.md")

	var g errgroup.Group

	var (
		settings     *claude.Settings
		mtSettings   time.Time
		localSet     *claude.Settings
		mtLocal      time.Time
		claudeMD     *claude.ClaudeMD
		mtClaudeMD   time.Time
		subClaudeMDs []claude.ClaudeMD
		mtSubMDs     map[string]time.Time
		gitInfo      *GitInfo
		memFiles     []claude.MemoryFile
	)

	g.Go(func() error {
		if info, err := os.Stat(settingsPath); err == nil {
			mtSettings = info.ModTime()
			settings, _ = claude.ParseSettings(settingsPath)
		}
		return nil
	})

	g.Go(func() error {
		if info, err := os.Stat(localSettingsPath); err == nil {
			mtLocal = info.ModTime()
			localSet, _ = claude.ParseSettings(localSettingsPath)
		}
		return nil
	})

	g.Go(func() error {
		if info, err := os.Stat(claudeMDPath); err == nil {
			mtClaudeMD = info.ModTime()
			claudeMD, _ = claude.ParseClaudeMD(claudeMDPath)
		}
		return nil
	})

	g.Go(func() error {
		mtSubMDs = make(map[string]time.Time)
		subClaudeMDs = scanSubClaudeMDs(l.ch.Path, claudeMDPath, mtSubMDs)
		return nil
	})

	g.Go(func() error {
		gitInfo = loadGitInfo(l.ch.Path)
		return nil
	})

	if l.claudeHome != "" {
		g.Go(func() error {
			memoryDir := filepath.Join(l.claudeHome, "projects", l.ch.ID, "memory")
			memFiles, _ = claude.ScanMemoryFiles(memoryDir)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	data.Settings = settings
	if !mtSettings.IsZero() {
		mtimes[settingsPath] = mtSettings
	}
	data.LocalSettings = localSet
	if !mtLocal.IsZero() {
		mtimes[localSettingsPath] = mtLocal
	}
	data.ClaudeMD = claudeMD
	if !mtClaudeMD.IsZero() {
		mtimes[claudeMDPath] = mtClaudeMD
	}

	data.SubClaudeMDs = subClaudeMDs
	for k, v := range mtSubMDs {
		mtimes[k] = v
	}

	data.GitInfo = gitInfo
	data.MemoryFiles = memFiles

	if data.Settings != nil {
		data.Hooks, _ = claude.ExtractHooks(data.Settings, "project")
		data.MCPServers, _ = claude.ExtractMCPServers(data.Settings, "project")
	}

	data.HealthIssues = claude.CheckHealth(&claude.HealthInput{
		ClaudeMD:   data.ClaudeMD,
		Settings:   data.Settings,
		MCPServers: data.MCPServers,
		Hooks:      data.Hooks,
	})

	return data, mtimes, nil
}

// LoadChannelData parses all data for a channel.
// claudeHome is ~/.claude (for global data like plugins/skills).
// Returns parsed data and file mtimes for caching.
func LoadChannelData(ch *Channel, claudeHome string) (*ChannelData, map[string]time.Time, error) {
	loader := &channelLoader{ch: ch, claudeHome: claudeHome}
	return loader.load()
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

const maxSubClaudeMDDepth = 4

func scanSubClaudeMDs(projectDir, rootClaudeMD string, mtimes map[string]time.Time) []claude.ClaudeMD {
	var result []claude.ClaudeMD
	baseDepth := strings.Count(projectDir, string(filepath.Separator))
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
			depth := strings.Count(path, string(filepath.Separator)) - baseDepth
			if depth > maxSubClaudeMDDepth {
				return filepath.SkipDir
			}
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

	var g errgroup.Group

	g.Go(func() error {
		if out, err := runGit(projectPath, "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
			info.Branch = strings.TrimSpace(out)
		}
		return nil
	})

	g.Go(func() error {
		if out, err := runGit(projectPath, "log", "-1", "--format=%h|%s|%ci"); err == nil {
			parts := strings.SplitN(strings.TrimSpace(out), "|", 3)
			if len(parts) == 3 {
				info.LastCommit = parts[0]
				info.LastCommitMsg = parts[1]
				info.LastCommitAt = parts[2]
			}
		}
		return nil
	})

	g.Go(func() error {
		if out, err := runGit(projectPath, "status", "--porcelain"); err == nil {
			trimmed := strings.TrimSpace(out)
			if trimmed == "" {
				info.DirtyFiles = 0
			} else {
				info.DirtyFiles = len(strings.Split(trimmed, "\n"))
			}
		}
		return nil
	})

	_ = g.Wait()

	return info
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return string(out), err
}
