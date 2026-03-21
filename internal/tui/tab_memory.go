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
		b.WriteString(section("Memory"))
		b.WriteString("    No memory files found\n")
		return b.String()
	}

	b.WriteString(section(fmt.Sprintf("Memory Files (%d)", len(ch.Data.MemoryFiles))))

	order, groups := orderedGroup(ch.Data.MemoryFiles, func(mf claude.MemoryFile) string {
		if mf.Type == "" {
			return "unknown"
		}
		return mf.Type
	})

	for _, t := range order {
		b.WriteString(fmt.Sprintf("\n    %s\n", headerStyle.Render("["+t+"]")))
		for _, mf := range groups[t] {
			b.WriteString(fmt.Sprintf("      %s\n", mf.Name))
			if mf.Description != "" {
				b.WriteString(fmt.Sprintf("        %s\n", labelStyle.Render(mf.Description)))
			}
		}
	}
	return b.String()
}
