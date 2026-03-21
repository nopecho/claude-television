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
		b.WriteString(section("Hooks"))
		b.WriteString("    No hooks registered\n")
		return b.String()
	}

	b.WriteString(section(fmt.Sprintf("Hooks (%d)", len(ch.Data.Hooks))))

	order, groups := orderedGroup(ch.Data.Hooks, func(h claude.HookDetail) string { return h.Event })
	for _, event := range order {
		b.WriteString(fmt.Sprintf("\n    %s\n", headerStyle.Render(event)))
		for _, h := range groups[event] {
			source := labelStyle.Render(fmt.Sprintf("[%s]", h.Source))
			b.WriteString(fmt.Sprintf("      %s %s\n", source, h.Type))
			if h.Matcher != "" {
				b.WriteString(fmt.Sprintf("        matcher: %s\n", h.Matcher))
			}
			if h.Command != "" {
				b.WriteString(fmt.Sprintf("        command: %s\n", h.Command))
			}
			if h.Async {
				b.WriteString("        async: true\n")
			}
			if h.Timeout > 0 {
				b.WriteString(fmt.Sprintf("        timeout: %ds\n", h.Timeout))
			}
		}
	}
	return b.String()
}
