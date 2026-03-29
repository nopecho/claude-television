package tui

import (
	"strings"
	"testing"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func modelWithTab(channels []channel.Channel, tab DetailTab) model {
	m := newModel(channels, defaultCfg())
	m.detailTab = tab
	return m
}

func TestRenderDetailContentString_NoChannel(t *testing.T) {
	m := newModel([]channel.Channel{}, defaultCfg())

	got := m.renderDetailContentString()
	if !strings.Contains(got, "No channel") {
		t.Errorf("with no channel: should show 'No channel', got: %q", got)
	}
}

func TestRenderDetailContentString_NilData(t *testing.T) {
	ch := channel.Channel{ID: "1", Name: "alpha", Path: "/tmp/alpha", Data: nil}
	m := newModel([]channel.Channel{ch}, defaultCfg())

	got := m.renderDetailContentString()
	if !strings.Contains(got, "No channel") {
		t.Errorf("with nil data: should show 'No channel', got: %q", got)
	}
}

// Settings tab
func TestRenderSettingsTab_Empty(t *testing.T) {
	ch := makeChannel("1", "alpha")
	m := modelWithTab([]channel.Channel{ch}, TabSettings)

	got := m.renderDetailContentString()
	if got == "" {
		t.Error("renderSettingsTab should not return empty string")
	}
	if !strings.Contains(got, "alpha") {
		t.Errorf("settings tab should show channel path, got: %q", got)
	}
}

func TestRenderSettingsTab_WithSettings(t *testing.T) {
	ch := makeChannel("1", "alpha")
	ch.Data.Settings = &claude.Settings{
		Model:    "claude-3-opus",
		Language: "en",
	}
	m := modelWithTab([]channel.Channel{ch}, TabSettings)

	got := m.renderDetailContentString()
	if !strings.Contains(got, "claude-3-opus") {
		t.Errorf("settings tab should show model, got: %q", got)
	}
	if !strings.Contains(got, "en") {
		t.Errorf("settings tab should show language, got: %q", got)
	}
}

func TestRenderSettingsTab_WithPermissions(t *testing.T) {
	ch := makeChannel("1", "alpha")
	ch.Data.Settings = &claude.Settings{
		Permissions: claude.Permissions{
			Allow: []string{"Bash"},
			Deny:  []string{"Write"},
		},
	}
	m := modelWithTab([]channel.Channel{ch}, TabSettings)

	got := m.renderDetailContentString()
	if !strings.Contains(got, "Bash") {
		t.Errorf("settings tab should show allowed permissions, got: %q", got)
	}
	if !strings.Contains(got, "Write") {
		t.Errorf("settings tab should show denied permissions, got: %q", got)
	}
}

func TestRenderSettingsTab_DiffWithGlobal(t *testing.T) {
	global := makeChannel("g", "global", withGlobal())
	global.Data.Settings = &claude.Settings{Model: "opus"}

	project := makeChannel("p", "project")
	project.Data.Settings = &claude.Settings{Model: "sonnet"}

	m := modelWithTab([]channel.Channel{global, project}, TabSettings)
	m.channelCursor = 1 // select project

	got := m.renderDetailContentString()
	if !strings.Contains(got, "sonnet") {
		t.Errorf("settings tab should show project model, got: %q", got)
	}
}

// ClaudeMD tab
func TestRenderClaudeMDTab_Empty(t *testing.T) {
	ch := makeChannel("1", "alpha")
	m := modelWithTab([]channel.Channel{ch}, TabClaudeMD)

	got := m.renderDetailContentString()
	if got == "" {
		t.Error("renderClaudeMDTab should not return empty string")
	}
}

func TestRenderClaudeMDTab_WithContent(t *testing.T) {
	ch := makeChannel("1", "alpha")
	ch.Data.ClaudeMD = &claude.ClaudeMD{
		Path:    "/project/CLAUDE.md",
		Content: "# Project Instructions\nDo something",
	}
	m := modelWithTab([]channel.Channel{ch}, TabClaudeMD)
	m.width = 120
	m.height = 40

	got := m.renderDetailContentString()
	if !strings.Contains(got, "Project") {
		t.Errorf("claudemd tab should show content, got: %q", got)
	}
}

// Hooks tab
func TestRenderHooksTab_Empty(t *testing.T) {
	ch := makeChannel("1", "alpha")
	m := modelWithTab([]channel.Channel{ch}, TabHooks)

	got := m.renderDetailContentString()
	if !strings.Contains(got, "No hooks") {
		t.Errorf("empty hooks tab should show 'No hooks', got: %q", got)
	}
}

func TestRenderHooksTab_WithHooks(t *testing.T) {
	ch := makeChannel("1", "alpha")
	ch.Data.Hooks = []claude.HookDetail{
		{Event: "PreToolUse", Type: "command", Command: "echo test", Source: "project"},
	}
	m := modelWithTab([]channel.Channel{ch}, TabHooks)

	got := m.renderDetailContentString()
	if !strings.Contains(got, "PreToolUse") {
		t.Errorf("hooks tab should show event name, got: %q", got)
	}
}

// MCP tab
func TestRenderMCPTab_Empty(t *testing.T) {
	ch := makeChannel("1", "alpha")
	m := modelWithTab([]channel.Channel{ch}, TabMCP)

	got := m.renderDetailContentString()
	if got == "" {
		t.Error("renderMCPTab should not return empty string")
	}
}

func TestRenderMCPTab_WithServers(t *testing.T) {
	ch := makeChannel("1", "alpha")
	ch.Data.MCPServers = []claude.MCPServer{
		{Name: "my-server", Type: "stdio"},
	}
	m := modelWithTab([]channel.Channel{ch}, TabMCP)

	got := m.renderDetailContentString()
	if !strings.Contains(got, "my-server") {
		t.Errorf("mcp tab should show server name, got: %q", got)
	}
}

// Plugins tab
func TestRenderPluginsTab_Empty(t *testing.T) {
	ch := makeChannel("1", "alpha")
	m := modelWithTab([]channel.Channel{ch}, TabPlugins)

	got := m.renderDetailContentString()
	if got == "" {
		t.Error("renderPluginsTab should not return empty string")
	}
}

// Health tab
func TestRenderHealthTab_Empty(t *testing.T) {
	ch := makeChannel("1", "alpha")
	m := modelWithTab([]channel.Channel{ch}, TabHealth)

	got := m.renderDetailContentString()
	if got == "" {
		t.Error("renderHealthTab should not return empty string")
	}
}

func TestRenderHealthTab_WithIssues(t *testing.T) {
	ch := makeChannel("1", "alpha")
	ch.Data.HealthIssues = []claude.HealthIssue{
		{Code: "MISSING_SETTINGS", Severity: claude.SeverityError, Message: "settings.json not found"},
	}
	m := modelWithTab([]channel.Channel{ch}, TabHealth)
	m.width = 120
	m.height = 40

	got := m.renderDetailContentString()
	if !strings.Contains(got, "settings.json") {
		t.Errorf("health tab should show issue message, got: %q", got)
	}
}

// Git tab
func TestRenderGitTab_NoData(t *testing.T) {
	ch := makeChannel("1", "alpha")
	m := modelWithTab([]channel.Channel{ch}, TabGit)

	got := m.renderDetailContentString()
	if got == "" {
		t.Error("renderGitTab should not return empty string")
	}
}

func TestRenderGitTab_WithData(t *testing.T) {
	ch := makeChannel("1", "alpha")
	ch.Data.GitInfo = &channel.GitInfo{
		Branch:        "main",
		LastCommit:    "abc1234",
		LastCommitMsg: "initial commit",
	}
	m := modelWithTab([]channel.Channel{ch}, TabGit)

	got := m.renderDetailContentString()
	if !strings.Contains(got, "main") {
		t.Errorf("git tab should show branch, got: %q", got)
	}
}

// Memory tab
func TestRenderMemoryTab_Empty(t *testing.T) {
	ch := makeChannel("1", "alpha")
	m := modelWithTab([]channel.Channel{ch}, TabMemory)

	got := m.renderDetailContentString()
	if got == "" {
		t.Error("renderMemoryTab should not return empty string")
	}
}

func TestRenderMemoryTab_WithFiles(t *testing.T) {
	ch := makeChannel("1", "alpha")
	ch.Data.MemoryFiles = []claude.MemoryFile{
		{Path: "/project/memory/notes.md", Name: "notes.md", Type: "manual"},
	}
	m := modelWithTab([]channel.Channel{ch}, TabMemory)

	got := m.renderDetailContentString()
	if !strings.Contains(got, "notes.md") {
		t.Errorf("memory tab should show file path, got: %q", got)
	}
}
