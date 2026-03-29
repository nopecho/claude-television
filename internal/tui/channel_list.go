package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderChannelList(height int) string {
	var b strings.Builder

	count := len(m.filtered)

	if count == 0 {
		// Empty state with friendly message
		b.WriteString("\n")
		b.WriteString(channelItemStyle.Render(emptyMsgStyle.Render("  No channels found.")))
		b.WriteString("\n\n")
		b.WriteString(channelItemStyle.Render(emptyHintStyle.Render("  Press Esc to clear")))
		return b.String()
	}

	type displayItem struct {
		isHeader bool
		header   string
		chIdx    int
	}

	var items []displayItem
	lastGroup := ""
	for _, idx := range m.filtered {
		ch := m.channels[idx]

		if ch.Group != "" && ch.Group != lastGroup {
			lastGroup = ch.Group
			items = append(items, displayItem{isHeader: true, header: ch.Group})
		} else if ch.Group == "" && lastGroup != "" {
			lastGroup = ""
		}

		items = append(items, displayItem{chIdx: idx})
	}

	cursorDisplayIdx := 0
	channelCount := 0
	for i, item := range items {
		if !item.isHeader {
			if channelCount == m.channelCursor {
				cursorDisplayIdx = i
				break
			}
			channelCount++
		}
	}

	// Selected items take 2 lines (name + path), so adjust visible height
	visibleHeight := height
	if visibleHeight < 1 {
		visibleHeight = 1
	}
	startIdx := cursorDisplayIdx - visibleHeight/2
	if startIdx < 0 {
		startIdx = 0
	}

	// Count display lines from startIdx to find how many items fit
	maxItems := startIdx + len(items) // upper bound
	displayLines := 0
	for i := startIdx; i < len(items); i++ {
		if !items[i].isHeader && i == cursorDisplayIdx {
			displayLines += 2 // selected item takes 2 lines
		} else {
			displayLines++
		}
		if displayLines >= visibleHeight {
			maxItems = i + 1
			break
		}
	}

	// Adjust start if we hit the end
	if maxItems > len(items) {
		maxItems = len(items)
	}

	listWidth := m.listWidth() - 2 // account for border padding

	channelIdx := 0
	for i, item := range items {
		if i < startIdx {
			if !item.isHeader {
				channelIdx++
			}
			continue
		}
		if i >= maxItems {
			break
		}

		if item.isHeader {
			// Group divider style: ── group ──────
			label := fmt.Sprintf(" %s ", item.header)
			labelWidth := lipgloss.Width(label)
			remaining := listWidth - labelWidth - 2
			if remaining < 2 {
				remaining = 2
			}
			divider := groupDividerStyle.Render("──") +
				groupHeaderStyle.Render(label) +
				groupDividerStyle.Render(strings.Repeat("─", remaining))
			b.WriteString(divider)
			b.WriteString("\n")
			continue
		}

		ch := m.channels[item.chIdx]
		prefix := "  "
		icon := statusIconStr(ch.Status)
		if ch.IsGlobal {
			prefix = "⚙ "
			icon = ""
		} else if ch.Pinned {
			prefix = pinIcon + " "
		}

		// Error badge
		badge := ""
		if ch.Data != nil && len(ch.Data.HealthIssues) > 0 {
			badge = " " + channelBadgeStyle.Render(fmt.Sprintf("⚠%d", len(ch.Data.HealthIssues)))
		}

		name := truncate(ch.Name, listWidth-8)
		line := fmt.Sprintf("%s%s %s%s", prefix, icon, name, badge)

		if channelIdx == m.channelCursor {
			// Selected: background highlight + bold + path on second line
			b.WriteString(channelSelectedStyle.
				Width(listWidth).
				Render("▸" + line))
			b.WriteString("\n")
			// Show path below selected item
			shortPath := truncate(ch.Path, listWidth-4)
			b.WriteString(channelSelectedPathStyle.
				Width(listWidth).
				Render("   " + shortPath))
			b.WriteString("\n")
		} else {
			b.WriteString(channelItemStyle.Render(" " + line))
			b.WriteString("\n")
		}
		channelIdx++
	}

	return b.String()
}
