package claude

import (
	"os/exec"
	"strings"
)

var deprecatedSettingKeys = []string{
	"apiKeyHelper",
	"disableAutoupdater",
	"autoUpdaterStatus",
}

type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
)

type HealthIssue struct {
	Code     string   `json:"code"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
}

type HealthInput struct {
	ClaudeMD   *ClaudeMD
	Settings   *Settings
	MCPServers []MCPServer
	Hooks      []HookDetail
	Plugins    []Plugin
	GitError   error
}

var dangerousPermissions = []string{
	"Bash(*)", "Edit(*)", "Write(*)",
}

func CheckHealth(input *HealthInput) []HealthIssue {
	var issues []HealthIssue

	issues = append(issues, checkGitError(input.GitError)...)
	if input.ClaudeMD != nil {
		issues = append(issues, checkClaudeMD(input.ClaudeMD)...)
	}
	if input.Settings != nil {
		issues = append(issues, checkPermissions(input.Settings)...)
		issues = append(issues, checkDeprecatedSettings(input.Settings)...)
	}
	issues = append(issues, checkMCPServers(input.MCPServers)...)
	issues = append(issues, checkMCPServerDuplicates(input.MCPServers)...)
	issues = append(issues, checkHooks(input.Hooks)...)
	issues = append(issues, checkPluginConflicts(input.Plugins)...)

	return issues
}

func checkClaudeMD(md *ClaudeMD) []HealthIssue {
	var issues []HealthIssue
	if md.LineCount > 200 {
		issues = append(issues, HealthIssue{
			Code:     "claudemd-too-long",
			Severity: SeverityWarning,
			Message:  "CLAUDE.md exceeds 200 lines — Claude Code may truncate it",
		})
	}
	return issues
}

func checkPermissions(s *Settings) []HealthIssue {
	var issues []HealthIssue
	for _, perm := range s.Permissions.Allow {
		for _, dangerous := range dangerousPermissions {
			if perm == dangerous {
				issues = append(issues, HealthIssue{
					Code:     "dangerous-permission",
					Severity: SeverityWarning,
					Message:  "Broad permission allowed: " + perm,
				})
			}
		}
	}
	return issues
}

func checkMCPServers(servers []MCPServer) []HealthIssue {
	var issues []HealthIssue
	for _, s := range servers {
		if s.Command == "" {
			continue
		}
		cmd := strings.Fields(s.Command)[0]
		if _, err := exec.LookPath(cmd); err != nil {
			issues = append(issues, HealthIssue{
				Code:     "mcp-command-missing",
				Severity: SeverityError,
				Message:  "MCP server '" + s.Name + "': command not found: " + cmd,
			})
		}
	}
	return issues
}

func checkHooks(hooks []HookDetail) []HealthIssue {
	var issues []HealthIssue
	for _, h := range hooks {
		if h.Command == "" || h.Type != "command" {
			continue
		}
		cmd := strings.Fields(h.Command)[0]
		if _, err := exec.LookPath(cmd); err != nil {
			issues = append(issues, HealthIssue{
				Code:     "hook-command-missing",
				Severity: SeverityError,
				Message:  "Hook '" + h.Event + "': command not found: " + cmd,
			})
		}
	}
	return issues
}

func checkGitError(err error) []HealthIssue {
	if err == nil {
		return nil
	}
	return []HealthIssue{{
		Code:     "git-error",
		Severity: SeverityError,
		Message:  "git command failed: " + err.Error(),
	}}
}

func checkPluginConflicts(plugins []Plugin) []HealthIssue {
	var issues []HealthIssue
	for _, p := range plugins {
		if p.Installed && !p.Enabled {
			issues = append(issues, HealthIssue{
				Code:     "plugin-not-enabled",
				Severity: SeverityWarning,
				Message:  "Plugin '" + p.Name + "' is installed but not enabled",
			})
		}
	}
	return issues
}

func checkDeprecatedSettings(s *Settings) []HealthIssue {
	var issues []HealthIssue
	for _, key := range deprecatedSettingKeys {
		if _, ok := s.Raw[key]; ok {
			issues = append(issues, HealthIssue{
				Code:     "deprecated-setting",
				Severity: SeverityWarning,
				Message:  "Deprecated setting in use: " + key,
			})
		}
	}
	return issues
}

func checkMCPServerDuplicates(servers []MCPServer) []HealthIssue {
	seen := make(map[string][]string)
	for _, s := range servers {
		seen[s.Name] = append(seen[s.Name], s.Source)
	}
	var issues []HealthIssue
	for name, sources := range seen {
		if len(sources) > 1 {
			issues = append(issues, HealthIssue{
				Code:     "mcp-server-duplicate",
				Severity: SeverityWarning,
				Message:  "MCP server '" + name + "' defined in multiple sources: " + strings.Join(sources, ", "),
			})
		}
	}
	return issues
}
