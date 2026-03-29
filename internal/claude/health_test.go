package claude

import (
	"errors"
	"strings"
	"testing"
)

func TestCheckHealth_ClaudeMDTooLong(t *testing.T) {
	md := &ClaudeMD{Path: "/test/CLAUDE.md", LineCount: 250, Content: strings.Repeat("line\n", 250)}
	issues := CheckHealth(&HealthInput{ClaudeMD: md})
	found := false
	for _, i := range issues {
		if i.Code == "claudemd-too-long" {
			found = true
		}
	}
	if !found {
		t.Error("expected claudemd-too-long issue")
	}
}

func TestCheckHealth_DangerousPermission(t *testing.T) {
	s := &Settings{Permissions: Permissions{Allow: []string{"Bash(*)", "Read"}}}
	issues := CheckHealth(&HealthInput{Settings: s})
	found := false
	for _, i := range issues {
		if i.Code == "dangerous-permission" {
			found = true
		}
	}
	if !found {
		t.Error("expected dangerous-permission issue")
	}
}

func TestCheckHealth_MissingMCPCommand(t *testing.T) {
	servers := []MCPServer{{Name: "test", Command: "/nonexistent/binary"}}
	issues := CheckHealth(&HealthInput{MCPServers: servers})
	found := false
	for _, i := range issues {
		if i.Code == "mcp-command-missing" {
			found = true
		}
	}
	if !found {
		t.Error("expected mcp-command-missing issue")
	}
}

func TestCheckHealth_MissingHookCommand(t *testing.T) {
	hooks := []HookDetail{{Event: "test", Command: "/nonexistent/hook", Type: "command"}}
	issues := CheckHealth(&HealthInput{Hooks: hooks})
	found := false
	for _, i := range issues {
		if i.Code == "hook-command-missing" {
			found = true
		}
	}
	if !found {
		t.Error("expected hook-command-missing issue")
	}
}

func TestCheckGitError(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		wantHit bool
	}{
		{"no error", nil, false},
		{"with error", errors.New("exit status 128"), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			issues := CheckHealth(&HealthInput{GitError: tc.err})
			found := false
			for _, i := range issues {
				if i.Code == "git-error" {
					found = true
				}
			}
			if found != tc.wantHit {
				t.Errorf("found=%v, want %v", found, tc.wantHit)
			}
		})
	}
}

func TestCheckPluginConflicts(t *testing.T) {
	cases := []struct {
		name    string
		plugins []Plugin
		wantHit bool
	}{
		{"no plugins", nil, false},
		{"installed and enabled", []Plugin{{Name: "foo", Installed: true, Enabled: true}}, false},
		{"installed but not enabled", []Plugin{{Name: "bar", Installed: true, Enabled: false}}, true},
		{"not installed not enabled", []Plugin{{Name: "baz", Installed: false, Enabled: false}}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			issues := CheckHealth(&HealthInput{Plugins: tc.plugins})
			found := false
			for _, i := range issues {
				if i.Code == "plugin-not-enabled" {
					found = true
				}
			}
			if found != tc.wantHit {
				t.Errorf("found=%v, want %v", found, tc.wantHit)
			}
		})
	}
}

func TestCheckDeprecatedSettings(t *testing.T) {
	cases := []struct {
		name    string
		raw     map[string]any
		wantHit bool
	}{
		{"no deprecated", map[string]any{"model": "claude-3"}, false},
		{"apiKeyHelper", map[string]any{"apiKeyHelper": "some-script"}, true},
		{"disableAutoupdater", map[string]any{"disableAutoupdater": true}, true},
		{"autoUpdaterStatus", map[string]any{"autoUpdaterStatus": "disabled"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := &Settings{Raw: tc.raw}
			issues := CheckHealth(&HealthInput{Settings: s})
			found := false
			for _, i := range issues {
				if i.Code == "deprecated-setting" {
					found = true
				}
			}
			if found != tc.wantHit {
				t.Errorf("found=%v, want %v", found, tc.wantHit)
			}
		})
	}
}

func TestCheckMCPServerDuplicates(t *testing.T) {
	cases := []struct {
		name    string
		servers []MCPServer
		wantHit bool
	}{
		{"no servers", nil, false},
		{"unique names", []MCPServer{{Name: "a", Source: "user"}, {Name: "b", Source: "project"}}, false},
		{"duplicate name different sources", []MCPServer{{Name: "a", Source: "user"}, {Name: "a", Source: "project"}}, true},
		{"single server", []MCPServer{{Name: "a", Source: "user"}}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			issues := CheckHealth(&HealthInput{MCPServers: tc.servers})
			found := false
			for _, i := range issues {
				if i.Code == "mcp-server-duplicate" {
					found = true
				}
			}
			if found != tc.wantHit {
				t.Errorf("found=%v, want %v", found, tc.wantHit)
			}
		})
	}
}

func TestCheckHealth_Healthy(t *testing.T) {
	md := &ClaudeMD{Path: "/test/CLAUDE.md", LineCount: 50, Content: "short"}
	s := &Settings{Permissions: Permissions{Allow: []string{"Read"}}}
	issues := CheckHealth(&HealthInput{ClaudeMD: md, Settings: s})
	for _, i := range issues {
		if i.Severity == SeverityError {
			t.Errorf("unexpected error issue: %s", i.Code)
		}
	}
}
