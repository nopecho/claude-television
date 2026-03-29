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
	runes := []rune(s)
	if max <= 1 || len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}

// card renders a card-style section with rounded border and inline title.
func card(title string, contentLines []string, width int) string {
	if width < 10 {
		width = 10
	}
	innerWidth := width - 4 // account for border + padding
	if innerWidth < 6 {
		innerWidth = 6
	}

	content := strings.Join(contentLines, "\n")

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(overlayColor).
		Padding(0, 1).
		Width(innerWidth)

	rendered := style.Render(content)

	if title != "" {
		rendered = injectBorderTitle(rendered, " "+cardTitleStyle.Render(title)+" ")
	}
	return rendered
}

// cardKV renders a key-value pair with aligned key width inside a card.
func cardKV(key, value string, keyWidth int) string {
	k := labelStyle.Render(fmt.Sprintf("%-*s", keyWidth, key))
	return fmt.Sprintf("%s  %s", k, valueStyle.Render(value))
}

// sectionHeader renders: ▎Section Title (legacy, used in cards)
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

// kv renders a key-value pair with aligned key width (legacy section style).
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
	b.WriteString("\n")
	b.WriteString("  " + emptyMsgStyle.Render(message) + "\n")
	if hint != "" {
		b.WriteString("  " + emptyHintStyle.Render(hint) + "\n")
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

// renderScrollbar renders a vertical mini scrollbar.
func renderScrollbar(viewportHeight, totalLines, offset int) string {
	if totalLines <= viewportHeight || viewportHeight < 3 {
		return strings.Repeat(" \n", viewportHeight)
	}

	trackHeight := viewportHeight
	thumbSize := max(1, trackHeight*viewportHeight/totalLines)
	thumbPos := 0
	if totalLines > viewportHeight {
		thumbPos = offset * (trackHeight - thumbSize) / (totalLines - viewportHeight)
	}
	if thumbPos+thumbSize > trackHeight {
		thumbPos = trackHeight - thumbSize
	}

	var sb strings.Builder
	for i := 0; i < trackHeight; i++ {
		if i >= thumbPos && i < thumbPos+thumbSize {
			sb.WriteString(scrollThumbStyle.Render("┃"))
		} else {
			sb.WriteString(scrollTrackStyle.Render("│"))
		}
		if i < trackHeight-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
