package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Name         string
	Path         string
	HasClaudeDir bool
	HasClaudeMD  bool
	HasSettings  bool
}

func ScanProjects(roots []string, ignore []string) ([]Project, error) {
	ignoreSet := make(map[string]bool, len(ignore))
	for _, ig := range ignore {
		ignoreSet[ig] = true
	}
	var projects []Project
	for _, root := range roots {
		expanded := expandHome(root)
		entries, err := os.ReadDir(expanded)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("scan root %s: %w", root, err)
		}
		for _, e := range entries {
			if !e.IsDir() || ignoreSet[e.Name()] || e.Name()[0] == '.' {
				continue
			}
			dirPath := filepath.Join(expanded, e.Name())
			p := Project{Name: e.Name(), Path: dirPath}
			if _, err := os.Stat(filepath.Join(dirPath, ".claude")); err == nil {
				p.HasClaudeDir = true
			}
			if _, err := os.Stat(filepath.Join(dirPath, "CLAUDE.md")); err == nil {
				p.HasClaudeMD = true
			}
			if _, err := os.Stat(filepath.Join(dirPath, ".claude", "settings.json")); err == nil {
				p.HasSettings = true
			}
			projects = append(projects, p)
		}
	}
	return projects, nil
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
