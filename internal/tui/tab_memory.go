package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderMemoryTab(ch *channel.Channel) string {
	var b strings.Builder

	if len(ch.Data.MemoryFiles) == 0 {
		return emptyState("Memory", "No memory files found", "Memory files are created by Claude automatically")
	}

	b.WriteString(section(fmt.Sprintf("Memory Files (%d)", len(ch.Data.MemoryFiles))))

	byType := map[string][]int{}
	var typeOrder []string
	for i, mf := range ch.Data.MemoryFiles {
		t := mf.Type
		if t == "" {
			t = "unknown"
		}
		if _, exists := byType[t]; !exists {
			typeOrder = append(typeOrder, t)
		}
		byType[t] = append(byType[t], i)
	}

	for _, t := range typeOrder {
		b.WriteString(sectionEmpty() + "\n")
		b.WriteString(sectionLine(sectionTitleStyle.Render("["+t+"]")) + "\n")
		for _, idx := range byType[t] {
			mf := ch.Data.MemoryFiles[idx]
			b.WriteString(sectionLine("  "+valueStyle.Render(mf.Name)) + "\n")
			if mf.Description != "" {
				b.WriteString(sectionLine("    "+labelStyle.Render(mf.Description)) + "\n")
			}
		}
	}
	return b.String()
}
