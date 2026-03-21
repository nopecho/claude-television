package tui

import (
	"fmt"
	"strings"
)

func (m model) renderHooksList() string {
	if len(m.data.Hooks) == 0 {
		return listItemStyle.Render("No hooks registered.")
	}
	var b strings.Builder
	for i, h := range m.data.Hooks {
		line := fmt.Sprintf("%s [%s]", h.Event, h.Source)
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("▸ " + line))
		} else {
			b.WriteString(listItemStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m model) renderHooksDetail() string {
	if len(m.data.Hooks) == 0 || m.cursor >= len(m.data.Hooks) {
		return ""
	}
	h := m.data.Hooks[m.cursor]
	var b strings.Builder
	b.WriteString(titleStyle.Render(h.Event) + "\n\n")
	b.WriteString(fmt.Sprintf("  event:   %s\n", h.Event))
	b.WriteString(fmt.Sprintf("  type:    %s\n", h.Type))
	if h.Matcher != "" {
		b.WriteString(fmt.Sprintf("  matcher: %s\n", h.Matcher))
	}
	b.WriteString(fmt.Sprintf("  command: %s\n", h.Command))
	if h.Async {
		b.WriteString(fmt.Sprintf("  async:   %v\n", h.Async))
	}
	if h.Timeout > 0 {
		b.WriteString(fmt.Sprintf("  timeout: %ds\n", h.Timeout))
	}
	b.WriteString(fmt.Sprintf("  source:  %s\n", h.Source))
	return b.String()
}
