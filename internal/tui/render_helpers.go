package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nopecho/claude-television/internal/channel"
)

func statusIconStr(s channel.ChannelStatus) string {
	switch s {
	case channel.StatusHealthy:
		return statusHealthy
	case channel.StatusWarning:
		return statusWarning
	case channel.StatusError:
		return statusError
	}
	return "?"
}

func boolIcon(b bool) string {
	if b {
		return statusHealthy
	}
	return statusError
}

func truncate(s string, max int) string {
	if max <= 1 || len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

// sectionHeader renders: ▎Section Title
func sectionHeader(title string) string {
	bar := sectionTitleStyle.Render("▎")
	text := sectionTitleStyle.Render(title)
	return bar + text
}

// sectionLine renders: │ content
func sectionLine(content string) string {
	bar := sectionBarStyle.Render("│")
	return bar + " " + content
}

// sectionEmpty renders: │
func sectionEmpty() string {
	return sectionBarStyle.Render("│")
}

// kv renders a key-value pair with aligned key width.
func kv(key, value string, keyWidth int) string {
	k := labelStyle.Render(fmt.Sprintf("%-*s", keyWidth, key))
	return sectionLine(fmt.Sprintf("%s  %s", k, valueStyle.Render(value)))
}

func section(title string) string {
	return "\n" + sectionHeader(title) + "\n"
}

func bullet(s string) string {
	return sectionLine("  • " + s)
}

func emptyState(title, message, hint string) string {
	var b strings.Builder
	b.WriteString(section(title))
	b.WriteString(sectionLine("  "+lipgloss.NewStyle().Foreground(subtextColor).Render(message)) + "\n")
	if hint != "" {
		b.WriteString(sectionLine("  "+lipgloss.NewStyle().Foreground(subtleColor).Render(hint)) + "\n")
	}
	return b.String()
}

func helpEntry(key, desc string) string {
	return helpKeyStyle.Render(key) + " " + helpDescStyle.Render(desc)
}

// orderedGroup groups items by a key function, preserving first-seen order.
func orderedGroup[T any](items []T, keyFn func(T) string) ([]string, map[string][]T) {
	groups := map[string][]T{}
	var order []string
	for _, item := range items {
		key := keyFn(item)
		if _, exists := groups[key]; !exists {
			order = append(order, key)
		}
		groups[key] = append(groups[key], item)
	}
	return order, groups
}
