package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderSettingsTab(ch *channel.Channel) string {
	var b strings.Builder

	b.WriteString(section("Project"))
	b.WriteString(kv("path", ch.Path, 8) + "\n")
	b.WriteString(kv("status", statusIconStr(ch.Status)+" "+string(ch.Status), 8) + "\n")

	if ch.Data.Settings != nil {
		b.WriteString(renderSettingsSection(ch.Data.Settings, "Project Settings"))
	}
	if ch.Data.LocalSettings != nil {
		b.WriteString(renderSettingsSection(ch.Data.LocalSettings, "Local Settings (override)"))
	}
	if ch.Data.Settings == nil && ch.Data.LocalSettings == nil {
		b.WriteString(emptyState("Settings", "No settings.json found", "Configure in .claude/settings.json"))
	}
	return b.String()
}

func renderSettingsSection(s *claude.Settings, title string) string {
	var b strings.Builder
	b.WriteString(section(title))

	if s.Model != "" {
		b.WriteString(kv("model", s.Model, 10) + "\n")
	}
	if s.Language != "" {
		b.WriteString(kv("language", s.Language, 10) + "\n")
	}
	if s.TeammateMode != "" {
		b.WriteString(kv("teammate", s.TeammateMode, 10) + "\n")
	}
	if s.PlansDirectory != "" {
		b.WriteString(kv("plans dir", s.PlansDirectory, 10) + "\n")
	}

	if len(s.Env) > 0 {
		b.WriteString(section("Environment"))
		for k, v := range s.Env {
			b.WriteString(sectionLine(fmt.Sprintf("  %s = %s", k, v)) + "\n")
		}
	}

	if len(s.Permissions.Allow) > 0 || len(s.Permissions.Deny) > 0 {
		b.WriteString(section("Permissions"))
		if len(s.Permissions.Allow) > 0 {
			b.WriteString(sectionLine(fmt.Sprintf("  Allow (%d):", len(s.Permissions.Allow))) + "\n")
			for _, p := range s.Permissions.Allow {
				b.WriteString(sectionLine(fmt.Sprintf("    %s %s", statusHealthy, p)) + "\n")
			}
		}
		if len(s.Permissions.Deny) > 0 {
			b.WriteString(sectionLine(fmt.Sprintf("  Deny (%d):", len(s.Permissions.Deny))) + "\n")
			for _, p := range s.Permissions.Deny {
				b.WriteString(sectionLine(fmt.Sprintf("    %s %s", statusError, p)) + "\n")
			}
		}
	}
	return b.String()
}
