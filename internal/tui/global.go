package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/claude"
)

var globalItems = []string{"Settings", "Local Settings", "CLAUDE.md"}

func (m model) renderGlobalList() string {
	var b strings.Builder
	for i, item := range globalItems {
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("▸ " + item))
		} else {
			b.WriteString(listItemStyle.Render("  " + item))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m model) renderGlobalDetail() string {
	switch m.cursor {
	case 0:
		return m.renderSettingsDetail(m.data.Settings, "settings.json")
	case 1:
		return m.renderSettingsDetail(m.data.LocalSettings, "settings.local.json")
	case 2:
		return m.renderClaudeMDDetail()
	}
	return ""
}

func (m model) renderSettingsDetail(s *claude.Settings, name string) string {
	if s == nil {
		return detailStyle.Render(name + " not found")
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render(name) + "\n\n")
	if s.Model != "" {
		b.WriteString(fmt.Sprintf("  model:     %s\n", s.Model))
	}
	if s.Language != "" {
		b.WriteString(fmt.Sprintf("  language:  %s\n", s.Language))
	}
	if s.TeammateMode != "" {
		b.WriteString(fmt.Sprintf("  teammate:  %s\n", s.TeammateMode))
	}
	if len(s.Env) > 0 {
		b.WriteString("\n  env:\n")
		for k, v := range s.Env {
			b.WriteString(fmt.Sprintf("    %s: %s\n", k, v))
		}
	}
	if len(s.Permissions.Allow) > 0 {
		b.WriteString(fmt.Sprintf("\n  permissions.allow: (%d rules)\n", len(s.Permissions.Allow)))
		for _, p := range s.Permissions.Allow {
			b.WriteString(fmt.Sprintf("    %s %s\n", statusIcon, p))
		}
	}
	if len(s.Permissions.Deny) > 0 {
		b.WriteString(fmt.Sprintf("\n  permissions.deny: (%d rules)\n", len(s.Permissions.Deny)))
		for _, p := range s.Permissions.Deny {
			b.WriteString(fmt.Sprintf("    %s %s\n", statusIconOff, p))
		}
	}
	return b.String()
}

func (m model) renderClaudeMDDetail() string {
	md := m.data.ClaudeMD
	if md == nil {
		return detailStyle.Render("CLAUDE.md not found")
	}
	var b strings.Builder
	b.WriteString(titleStyle.Render("CLAUDE.md") + "\n\n")
	b.WriteString(fmt.Sprintf("  Lines: %d\n\n", md.LineCount))
	if len(md.Sections) > 0 {
		b.WriteString("  Sections:\n")
		for _, s := range md.Sections {
			b.WriteString(fmt.Sprintf("    • %s\n", s))
		}
	}
	return b.String()
}
