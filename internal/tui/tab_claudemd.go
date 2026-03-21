package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderClaudeMDTab(ch *channel.Channel) string {
	var b strings.Builder

	if ch.Data.ClaudeMD != nil {
		md := ch.Data.ClaudeMD
		b.WriteString(section("CLAUDE.md"))
		b.WriteString(kv("path", md.Path, 6) + "\n")
		b.WriteString(kv("lines", fmt.Sprintf("%d", md.LineCount), 6) + "\n")

		if len(md.Sections) > 0 {
			b.WriteString(section("Sections"))
			for _, s := range md.Sections {
				b.WriteString(bullet(s) + "\n")
			}
		}

		b.WriteString(section("Preview"))
		lines := strings.Split(md.Content, "\n")
		max := 20
		if len(lines) < max {
			max = len(lines)
		}
		for _, l := range lines[:max] {
			b.WriteString(sectionLine("  "+l) + "\n")
		}
		if len(lines) > 20 {
			b.WriteString(sectionLine(lipgloss.NewStyle().Foreground(subtextColor).Render(fmt.Sprintf("  ... (%d more lines)", len(lines)-20))) + "\n")
		}
	} else {
		b.WriteString(emptyState("CLAUDE.md", "Not found", "Create CLAUDE.md in project root"))
	}

	if len(ch.Data.SubClaudeMDs) > 0 {
		b.WriteString(section(fmt.Sprintf("Sub CLAUDE.md files (%d)", len(ch.Data.SubClaudeMDs))))
		for _, md := range ch.Data.SubClaudeMDs {
			b.WriteString(sectionLine(fmt.Sprintf("  %s (%d lines)", md.Path, md.LineCount)) + "\n")
			for _, s := range md.Sections {
				b.WriteString(sectionLine("    • "+s) + "\n")
			}
		}
	}
	return b.String()
}
