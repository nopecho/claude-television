package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderMemoryTab(ch *channel.Channel) string {
	var b strings.Builder

	if len(ch.Data.MemoryFiles) == 0 {
		return emptyState("Memory", "No memory files found", "Memory files are created by Claude automatically")
	}

	b.WriteString(section(fmt.Sprintf("Memory Files (%d)", len(ch.Data.MemoryFiles))))

	order, groups := orderedGroup(ch.Data.MemoryFiles, func(mf claude.MemoryFile) string {
		if mf.Type == "" {
			return "unknown"
		}
		return mf.Type
	})

	for _, t := range order {
		b.WriteString(sectionEmpty() + "\n")
		b.WriteString(sectionLine(sectionTitleStyle.Render("["+t+"]")) + "\n")
		for _, mf := range groups[t] {
			b.WriteString(sectionLine("  "+valueStyle.Render(mf.Name)) + "\n")
			if mf.Description != "" {
				b.WriteString(sectionLine("    "+labelStyle.Render(mf.Description)) + "\n")
			}
		}
	}
	return b.String()
}
