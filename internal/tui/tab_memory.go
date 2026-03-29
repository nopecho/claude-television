package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderMemoryTab(ch *channel.Channel) string {
	var b strings.Builder
	w := m.detailContentWidth()

	if len(ch.Data.MemoryFiles) == 0 {
		b.WriteString(card("Memory", []string{
			emptyMsgStyle.Render("No memory files found"),
			emptyHintStyle.Render("Memory files are created by Claude automatically"),
		}, w))
		return b.String()
	}

	order, groups := orderedGroup(ch.Data.MemoryFiles, func(mf claude.MemoryFile) string {
		if mf.Type == "" {
			return "unknown"
		}
		return mf.Type
	})

	for _, t := range order {
		var lines []string
		for _, mf := range groups[t] {
			lines = append(lines, valueStyle.Render(mf.Name))
			if mf.Description != "" {
				lines = append(lines, "  "+labelStyle.Render(mf.Description))
			}
		}
		b.WriteString(card(fmt.Sprintf("%s (%d)", t, len(groups[t])), lines, w))
		b.WriteString("\n")
	}
	return b.String()
}
