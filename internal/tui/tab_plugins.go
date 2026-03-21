package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderPluginsTab(ch *channel.Channel) string {
	var b strings.Builder

	if len(ch.Data.Plugins) == 0 {
		b.WriteString(emptyState("Plugins", "No plugins installed", "Install plugins via Claude Code marketplace"))
	} else {
		b.WriteString(section(fmt.Sprintf("Plugins (%d)", len(ch.Data.Plugins))))
		for _, p := range ch.Data.Plugins {
			icon := boolIcon(p.Enabled)
			b.WriteString(sectionEmpty() + "\n")
			b.WriteString(sectionLine(fmt.Sprintf("  %s %s", icon, sectionTitleStyle.Render(p.Name))) + "\n")
			if p.Marketplace != "" {
				b.WriteString(kv("marketplace", p.Marketplace, 12) + "\n")
			}
			if p.Version != "" {
				b.WriteString(kv("version", p.Version, 12) + "\n")
			}
			if p.InstallPath != "" {
				b.WriteString(kv("path", p.InstallPath, 12) + "\n")
			}
			if !p.Installed {
				b.WriteString(sectionLine("    "+labelStyle.Render("(not installed)")) + "\n")
			}
		}
	}

	if len(ch.Data.LocalSkills) == 0 {
		b.WriteString(emptyState("Skills", "No local skills found", "Skills are located in .gemini/skills/"))
	} else {
		b.WriteString(section(fmt.Sprintf("Skills (%d)", len(ch.Data.LocalSkills))))
		for _, s := range ch.Data.LocalSkills {
			b.WriteString(sectionLine(fmt.Sprintf("  %s %s", statusHealthy, valueStyle.Render(s.Name))) + "\n")
			b.WriteString(kv("path", s.Path, 6) + "\n")
		}
	}

	return b.String()
}
