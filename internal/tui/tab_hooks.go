package tui

import (
	"fmt"
	"strings"

	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/claude"
)

func (m model) renderHooksTab(ch *channel.Channel) string {
	var b strings.Builder
	w := m.detailContentWidth()

	if len(ch.Data.Hooks) == 0 {
		b.WriteString(card("Hooks", []string{
			emptyMsgStyle.Render("No hooks registered"),
			emptyHintStyle.Render("Configure hooks in .claude/settings.json"),
		}, w))
		return b.String()
	}

	order, groups := orderedGroup(ch.Data.Hooks, func(h claude.HookDetail) string { return h.Event })
	for _, event := range order {
		var lines []string
		for _, h := range groups[event] {
			source := labelStyle.Render(fmt.Sprintf("[%s]", h.Source))
			lines = append(lines, fmt.Sprintf("%s %s", source, h.Type))
			if h.Matcher != "" {
				lines = append(lines, fmt.Sprintf("  matcher: %s", h.Matcher))
			}
			if h.Command != "" {
				lines = append(lines, fmt.Sprintf("  command: %s", h.Command))
			}
			if h.Async {
				lines = append(lines, "  async: true")
			}
			if h.Timeout > 0 {
				lines = append(lines, fmt.Sprintf("  timeout: %ds", h.Timeout))
			}
			lines = append(lines, "")
		}
		b.WriteString(card(event, lines, w))
		b.WriteString("\n")
	}
	return b.String()
}
