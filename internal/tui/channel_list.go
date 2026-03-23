package tui

import (
	"fmt"
	"strings"
)

func (m model) renderChannelList(height int) string {
	var b strings.Builder

	count := len(m.filtered)

	if count == 0 {
		b.WriteString(channelItemStyle.Render("  No channels found."))
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

	visibleHeight := height
	if visibleHeight < 1 {
		visibleHeight = 1
	}
	startIdx := cursorDisplayIdx - visibleHeight/2
	if startIdx < 0 {
		startIdx = 0
	}
	if startIdx+visibleHeight > len(items) {
		startIdx = len(items) - visibleHeight
		if startIdx < 0 {
			startIdx = 0
		}
	}

	channelIdx := 0
	for i, item := range items {
		if i < startIdx {
			if !item.isHeader {
				channelIdx++
			}
			continue
		}
		if i >= startIdx+visibleHeight {
			break
		}

		if item.isHeader {
			b.WriteString(groupHeaderStyle.Render(fmt.Sprintf(" [%s]", item.header)))
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

		name := truncate(ch.Name, 22)
		line := fmt.Sprintf("%s%s %s", prefix, icon, name)

		if channelIdx == m.channelCursor {
			b.WriteString(channelSelectedStyle.Render("▸" + line))
		} else {
			b.WriteString(channelItemStyle.Render(" " + line))
		}
		b.WriteString("\n")
		channelIdx++
	}

	return b.String()
}
