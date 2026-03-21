package claude

import (
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
