package claude_test

import (
	"path/filepath"
	"testing"

	"github.com/nopecho/claude-television/internal/claude"
)

func TestParseInstalledPlugins(t *testing.T) {
	path := filepath.Join("testdata", "installed_plugins.json")
	plugins, err := claude.ParseInstalledPlugins(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, ok := plugins["superpowers@claude-plugins-official"]
	if !ok {
		t.Fatal("expected superpowers plugin to exist")
	}
	if p.Version != "5.0.5" {
		t.Errorf("version = %q, want %q", p.Version, "5.0.5")
	}
	if p.Scope != "global" {
		t.Errorf("scope = %q, want %q", p.Scope, "global")
	}
	if p.GitCommitSha != "abc123" {
		t.Errorf("gitCommitSha = %q, want %q", p.GitCommitSha, "abc123")
	}
}

func TestMergePluginData(t *testing.T) {
	installed := map[string]claude.InstalledPlugin{
		"superpowers@claude-plugins-official": {
			Version:     "5.0.5",
			InstallPath: "/some/path",
		},
	}
	enabled := map[string]bool{
		"superpowers@claude-plugins-official": true,
		"obsidian@obsidian-skills":            false,
	}

	result := claude.MergePluginData(installed, enabled)
	if len(result) != 2 {
		t.Fatalf("result len = %d, want 2", len(result))
	}

	byKey := make(map[string]claude.Plugin)
	for _, p := range result {
		byKey[p.Key] = p
	}

	sp := byKey["superpowers@claude-plugins-official"]
	if !sp.Installed {
		t.Error("superpowers should be installed")
	}
	if !sp.Enabled {
		t.Error("superpowers should be enabled")
	}
	if sp.Name != "superpowers" {
		t.Errorf("name = %q, want %q", sp.Name, "superpowers")
	}
	if sp.Marketplace != "claude-plugins-official" {
		t.Errorf("marketplace = %q, want %q", sp.Marketplace, "claude-plugins-official")
	}
	if sp.Version != "5.0.5" {
		t.Errorf("version = %q, want %q", sp.Version, "5.0.5")
	}

	ob := byKey["obsidian@obsidian-skills"]
	if ob.Installed {
		t.Error("obsidian should not be installed")
	}
	if ob.Enabled {
		t.Error("obsidian should not be enabled")
	}
}
