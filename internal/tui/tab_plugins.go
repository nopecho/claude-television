package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderPluginsTab(ch *channel.Channel) string {
	var b strings.Builder
	w := m.detailContentWidth()

	if len(ch.Data.Plugins) == 0 {
		b.WriteString(card("Plugins", []string{
			emptyMsgStyle.Render("No plugins installed"),
			emptyHintStyle.Render("Install plugins via Claude Code marketplace"),
		}, w))
		b.WriteString("\n")
	} else {
		for _, p := range ch.Data.Plugins {
			icon := boolIcon(p.Enabled)
			var lines []string
			lines = append(lines, fmt.Sprintf("%s %s", icon, labelStyle.Render(func() string {
				if p.Enabled {
					return "enabled"
				}
				return "disabled"
			}())))
			if p.Marketplace != "" {
				lines = append(lines, cardKV("marketplace", p.Marketplace, 12))
			}
			if p.Version != "" {
				lines = append(lines, cardKV("version", p.Version, 12))
			}
			if p.InstallPath != "" {
				lines = append(lines, cardKV("path", p.InstallPath, 12))
			}
			if !p.Installed {
				lines = append(lines, labelStyle.Render("(not installed)"))
			}
			b.WriteString(card(p.Name, lines, w))
			b.WriteString("\n")
		}
	}

	if len(ch.Data.LocalSkills) == 0 {
		b.WriteString(card("Skills", []string{
			emptyMsgStyle.Render("No local skills found"),
			emptyHintStyle.Render("Skills are located in .gemini/skills/"),
		}, w))
	} else {
		var lines []string
		for _, s := range ch.Data.LocalSkills {
			lines = append(lines, fmt.Sprintf("%s %s", statusHealthy, valueStyle.Render(s.Name)))
			lines = append(lines, fmt.Sprintf("  %s", labelStyle.Render(s.Path)))
		}
		b.WriteString(card(fmt.Sprintf("Skills (%d)", len(ch.Data.LocalSkills)), lines, w))
	}

	return b.String()
}
