package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type InstalledPlugin struct {
	Version      string `json:"version"`
	Scope        string `json:"scope"`
	InstallPath  string `json:"installPath"`
	InstalledAt  string `json:"installedAt"`
	GitCommitSha string `json:"gitCommitSha"`
}

type Plugin struct {
	Key         string
	Name        string
	Marketplace string
	Version     string
	Enabled     bool
	Installed   bool
	InstallPath string
}

func ParseInstalledPlugins(path string) (map[string]InstalledPlugin, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read installed_plugins: %w", err)
	}
	var plugins map[string]InstalledPlugin
	if err := json.Unmarshal(data, &plugins); err != nil {
		return nil, fmt.Errorf("parse installed_plugins: %w", err)
	}
	return plugins, nil
}

func MergePluginData(installed map[string]InstalledPlugin, enabled map[string]bool) []Plugin {
	seen := make(map[string]bool)
	result := []Plugin{}
	for key, ip := range installed {
		name, marketplace := splitPluginKey(key)
		result = append(result, Plugin{
			Key: key, Name: name, Marketplace: marketplace,
			Version: ip.Version, Enabled: enabled[key],
			Installed: true, InstallPath: ip.InstallPath,
		})
		seen[key] = true
	}
	for key, en := range enabled {
		if seen[key] {
			continue
		}
		name, marketplace := splitPluginKey(key)
		result = append(result, Plugin{
			Key: key, Name: name, Marketplace: marketplace,
			Enabled: en, Installed: false,
		})
	}
	return result
}

func splitPluginKey(key string) (name, marketplace string) {
	parts := strings.SplitN(key, "@", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return key, ""
}
