package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderHooksTab(ch *channel.Channel) string {
	var b strings.Builder

	if len(ch.Data.Hooks) == 0 {
		return emptyState("Hooks", "No hooks registered", "Configure hooks in .claude/settings.json")
	}

	b.WriteString(section(fmt.Sprintf("Hooks (%d)", len(ch.Data.Hooks))))

	order, groups := orderedGroup(ch.Data.Hooks, func(h claude.HookDetail) string { return h.Event })
	for _, event := range order {
		b.WriteString(sectionEmpty() + "\n")
		b.WriteString(sectionLine(sectionTitleStyle.Render(event)) + "\n")
		for _, h := range groups[event] {
			source := labelStyle.Render(fmt.Sprintf("[%s]", h.Source))
			b.WriteString(sectionLine(fmt.Sprintf("  %s %s", source, h.Type)) + "\n")
			if h.Matcher != "" {
				b.WriteString(sectionLine(fmt.Sprintf("    matcher: %s", h.Matcher)) + "\n")
			}
			if h.Command != "" {
				b.WriteString(sectionLine(fmt.Sprintf("    command: %s", h.Command)) + "\n")
			}
			if h.Async {
				b.WriteString(sectionLine("    async: true") + "\n")
			}
			if h.Timeout > 0 {
				b.WriteString(sectionLine(fmt.Sprintf("    timeout: %ds", h.Timeout)) + "\n")
			}
		}
	}
	return b.String()
}
