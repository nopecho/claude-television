package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderPluginsTab(ch *channel.Channel) string {
	var b strings.Builder

	if len(ch.Data.Plugins) == 0 {
		b.WriteString(section("Plugins"))
		b.WriteString("    No plugins installed\n")
	} else {
		b.WriteString(section(fmt.Sprintf("Plugins (%d)", len(ch.Data.Plugins))))
		for _, p := range ch.Data.Plugins {
			icon := boolIcon(p.Enabled)
			b.WriteString(fmt.Sprintf("\n    %s %s\n", icon, headerStyle.Render(p.Name)))
			if p.Marketplace != "" {
				b.WriteString(kv("      marketplace", p.Marketplace) + "\n")
			}
			if p.Version != "" {
				b.WriteString(kv("      version", p.Version) + "\n")
			}
			if p.InstallPath != "" {
				b.WriteString(kv("      path", p.InstallPath) + "\n")
			}
			if !p.Installed {
				b.WriteString("      " + labelStyle.Render("(not installed)") + "\n")
			}
		}
	}

	b.WriteString("\n")
	if len(ch.Data.LocalSkills) == 0 {
		b.WriteString(section("Skills"))
		b.WriteString("    No local skills found\n")
	} else {
		b.WriteString(section(fmt.Sprintf("Skills (%d)", len(ch.Data.LocalSkills))))
		for _, s := range ch.Data.LocalSkills {
			b.WriteString(fmt.Sprintf("    %s %s\n", statusHealthy, s.Name))
			b.WriteString(kv("      path", s.Path) + "\n")
		}
	}

	return b.String()
}
