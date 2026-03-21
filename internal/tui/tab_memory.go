package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) renderMemoryTab(ch *channel.Channel) string {
	var b strings.Builder

	if len(ch.Data.MemoryFiles) == 0 {
		b.WriteString(section("Memory"))
		b.WriteString("    No memory files found\n")
		return b.String()
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
		b.WriteString(fmt.Sprintf("\n    %s\n", headerStyle.Render("["+t+"]")))
		for _, idx := range byType[t] {
			mf := ch.Data.MemoryFiles[idx]
			b.WriteString(fmt.Sprintf("      %s\n", mf.Name))
			if mf.Description != "" {
				b.WriteString(fmt.Sprintf("        %s\n", labelStyle.Render(mf.Description)))
			}
		}
	}
	return b.String()
}
