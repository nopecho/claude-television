package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) View() string {
	if m.width == 0 {
		return "  Loading..."
	}

	listWidth := m.listWidth()
	detailWidth := m.width - listWidth - 4
	if detailWidth < 10 {
		detailWidth = 10
	}
	contentHeight := m.height - 5 // header + help + borders
	if contentHeight < 1 {
		contentHeight = 1
	}

	header := m.renderHeader()

	// Channel list panel
	listContent := m.renderChannelList(contentHeight)
	listTitle := fmt.Sprintf(" Channels (%d) ", len(m.filtered))

	var listBorder lipgloss.Style
	var listTitleStyle lipgloss.Style
	if m.focus == listPanel {
		listBorder = focusedBorderStyle
		listTitleStyle = focusedTitleStyle
	} else {
		listBorder = unfocusedBorderStyle
		listTitleStyle = unfocusedTitleStyle
	}
	listBox := listBorder.
		Width(listWidth).
		Height(contentHeight).
		BorderTop(true).
		Render(listContent)
	listBox = injectBorderTitle(listBox, listTitleStyle.Render(listTitle))

	// Detail panel
	tabBar := m.renderDetailTabs()
	vpView := m.viewport.View()

	// Scroll indicator in footer
	scrollInfo := ""
	if m.viewport.TotalLineCount() > m.viewport.VisibleLineCount() {
		pct := int(m.viewport.ScrollPercent() * 100)
		scrollInfo = fmt.Sprintf(" %d%% ", pct)
	}

	detailContent := tabBar + "\n" + vpView

	ch := m.selectedChannel()
	detailTitle := ""
	if ch != nil {
		detailTitle = fmt.Sprintf(" %s ", ch.Name)
	}

	var detailBorder lipgloss.Style
	var detailTitleStyle lipgloss.Style
	if m.focus == detailPanel {
		detailBorder = focusedBorderStyle
		detailTitleStyle = focusedTitleStyle
	} else {
		detailBorder = unfocusedBorderStyle
		detailTitleStyle = unfocusedTitleStyle
	}
	detailBox := detailBorder.
		Width(detailWidth).
		Height(contentHeight).
		Render(detailContent)
	detailBox = injectBorderTitle(detailBox, detailTitleStyle.Render(detailTitle))
	if scrollInfo != "" {
		detailBox = injectBorderFooter(detailBox, labelStyle.Render(scrollInfo))
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, listBox, detailBox)

	help := m.renderHelpBar()

	return header + content + "\n" + help
}

// renderHeader renders the unified header bar.
func (m model) renderHeader() string {
	if m.searching || m.contentSearching || m.grouping {
		left := headerAppStyle.Render(" ctv")
		search := "  " + m.searchInput.View()
		if m.contentSearching {
			if m.searchInput.Value() == "" {
				search += " " + labelStyle.Render("(content search)")
			} else if len(m.contentMatches) == 0 {
				search += " " + labelStyle.Render("(no matches)")
			} else {
				search += " " + labelStyle.Render(fmt.Sprintf("(%d/%d matches)", m.contentMatchIdx+1, len(m.contentMatches)))
			}
		} else if m.grouping {
			search += " " + labelStyle.Render("(enter to set group, empty to clear)")
		} else {
			search += " " + labelStyle.Render(fmt.Sprintf("(%d matches)", len(m.filtered)))
		}
		return left + search + "\n"
	}

	sep := headerSepStyle.Render(" ── ")

	left := headerAppStyle.Render(" ctv")

	// Selected channel name
	ch := m.selectedChannel()
	channelName := ""
	if ch != nil {
		channelName = headerChannelStyle.Render(ch.Name)
	}

	// Status badges
	healthy, warning, errCount := m.countStatuses()
	var badges []string
	if healthy > 0 {
		badges = append(badges, headerBadgeHealthy.Render(fmt.Sprintf("● %d", healthy)))
	}
	if warning > 0 {
		badges = append(badges, headerBadgeWarning.Render(fmt.Sprintf("○ %d", warning)))
	}
	if errCount > 0 {
		badges = append(badges, headerBadgeError.Render(fmt.Sprintf("✕ %d", errCount)))
	}

	total := len(m.channels)
	for _, ch := range m.channels {
		if ch.IsGlobal {
			total--
		}
	}
	countStr := headerCountStyle.Render(fmt.Sprintf("%d channels", total))

	header := left + sep + channelName + sep + strings.Join(badges, "  ") + sep + countStr
	return header + "\n"
}

func (m model) countStatuses() (healthy, warning, errCount int) {
	for _, ch := range m.channels {
		if ch.IsGlobal {
			continue
		}
		switch ch.Status {
		case channel.StatusHealthy:
			healthy++
		case channel.StatusWarning:
			warning++
		case channel.StatusError:
			errCount++
		}
	}
	return
}

func (m model) renderDetailTabs() string {
	var tabs []string

	ch := m.selectedChannel()

	for i, name := range detailTabNames {
		badge := m.tabBadge(ch, DetailTab(i))
		label := fmt.Sprintf("%d %s", i+1, name)
		if badge != "" {
			label += " " + tabBadgeStyle.Render(badge)
		}
		if DetailTab(i) == m.detailTab {
			tabs = append(tabs, activeTabStyle.Render(label))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(label))
		}
	}

	return " " + strings.Join(tabs, "")
}

// tabBadge returns a badge string for tabs with issues.
func (m model) tabBadge(ch *channel.Channel, tab DetailTab) string {
	if ch == nil || ch.Data == nil {
		return ""
	}
	switch tab {
	case TabHealth:
		count := len(ch.Data.HealthIssues)
		if count > 0 {
			return fmt.Sprintf("%d", count)
		}
	case TabGit:
		if ch.Data.GitInfo != nil && ch.Data.GitInfo.DirtyFiles > 0 {
			return "●"
		}
	}
	return ""
}

func (m model) renderHelpBar() string {
	var entries []string

	if m.searching || m.contentSearching || m.grouping {
		entries = []string{
			helpEntry("Enter", "confirm"),
			helpEntry("Esc", "cancel"),
		}
		if m.contentSearching && len(m.contentMatches) > 0 {
			entries = append(entries, helpEntry("n/N", "next/prev"))
		}
	} else if m.focus == detailPanel {
		entries = []string{
			helpEntry("j/k", "scroll"),
			helpEntry("1-8", "tabs"),
			helpEntry("Tab", "list"),
			helpEntry("/", "search"),
			helpEntry("?", "content"),
			helpEntry("e", "edit"),
			helpEntry("q", "quit"),
		}
	} else {
		entries = []string{
			helpEntry("j/k", "move"),
			helpEntry("1-8", "tabs"),
			helpEntry("l/Tab", "detail"),
			helpEntry("/", "search"),
			helpEntry("?", "content"),
			helpEntry("g", "group"),
			helpEntry("p", "pin"),
			helpEntry("e", "edit"),
			helpEntry("q", "quit"),
		}
	}

	helpText := "  " + strings.Join(entries, "  ")

	// Add scroll position on the right
	if !m.searching && !m.contentSearching && !m.grouping {
		if m.viewport.TotalLineCount() > m.viewport.VisibleLineCount() {
			pct := int(m.viewport.ScrollPercent() * 100)
			scrollStr := helpScrollStyle.Render(fmt.Sprintf("%d%%", pct))
			padding := m.width - lipgloss.Width(helpText) - lipgloss.Width(scrollStr) - 2
			if padding > 0 {
				helpText += strings.Repeat(" ", padding) + scrollStr
			}
		}
	}

	return helpText
}

// borderCharWidth is the byte length of the "─" UTF-8 character used in borders.
var borderCharWidth = len("─")

// injectBorderTitle places a styled title string into the top border line.
func injectBorderTitle(box, title string) string {
	if title == "" {
		return box
	}
	lines := strings.Split(box, "\n")
	if len(lines) == 0 {
		return box
	}
	first := lines[0]
	cornerEnd := strings.Index(first, "─")
	if cornerEnd < 0 {
		return box
	}
	insertAt := cornerEnd + borderCharWidth
	visualWidth := lipgloss.Width(title)
	replaceBytes := visualWidth * borderCharWidth
	if insertAt+replaceBytes > len(first) {
		return box
	}
	lines[0] = first[:insertAt] + title + first[insertAt+replaceBytes:]
	return strings.Join(lines, "\n")
}

// injectBorderFooter places a styled string into the bottom border line.
func injectBorderFooter(box, footer string) string {
	if footer == "" {
		return box
	}
	lines := strings.Split(box, "\n")
	if len(lines) < 2 {
		return box
	}
	last := lines[len(lines)-1]
	cornerEnd := strings.Index(last, "─")
	if cornerEnd < 0 {
		return box
	}
	insertAt := cornerEnd + borderCharWidth
	visualWidth := lipgloss.Width(footer)
	replaceBytes := visualWidth * borderCharWidth
	if insertAt+replaceBytes > len(last) {
		return box
	}
	lines[len(lines)-1] = last[:insertAt] + footer + last[insertAt+replaceBytes:]
	return strings.Join(lines, "\n")
}
