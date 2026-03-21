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
	contentHeight := m.height - 7 // header + tab bar + help + borders

	header := titleStyle.Render("ctv")
	if m.searching || m.contentSearching || m.grouping {
		header += "  " + m.searchInput.View()
		if m.contentSearching {
			header += " " + labelStyle.Render("(content search)")
		} else if m.grouping {
			header += " " + labelStyle.Render("(enter to set group, empty to clear)")
		} else {
			header += " " + labelStyle.Render(fmt.Sprintf("(%d matches)", len(m.filtered)))
		}
	} else {
		header += "  " + m.renderSummary()
	}
	header += "\n"

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

func (m model) renderSummary() string {
	healthy, warning, errCount, issues := 0, 0, 0, 0
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
		if ch.Data != nil {
			issues += len(ch.Data.HealthIssues)
		}
	}
	summary := labelStyle.Render(fmt.Sprintf("%d channels", healthy+warning+errCount))
	if errCount > 0 {
		summary += " " + lipgloss.NewStyle().Foreground(errorColor).Render(fmt.Sprintf("%d err", errCount))
	}
	if issues > 0 {
		summary += " " + lipgloss.NewStyle().Foreground(warningColor).Render(fmt.Sprintf("%d issues", issues))
	}
	return summary
}

func (m model) renderDetailTabs() string {
	var tabs []string
	var underlines []string

	for i, name := range detailTabNames {
		if DetailTab(i) == m.detailTab {
			tabs = append(tabs, activeTabStyle.Render(name))
			underlines = append(underlines, tabUnderlineStyle.Render(strings.Repeat("━", len(name))))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(name))
			underlines = append(underlines, strings.Repeat(" ", len(name)+2))
		}
	}

	return " " + strings.Join(tabs, "") + "\n " + strings.Join(underlines, "")
}

func (m model) renderHelpBar() string {
	var entries []string

	if m.searching || m.contentSearching || m.grouping {
		entries = []string{
			helpEntry("type", "input"),
			helpEntry("Enter", "confirm"),
			helpEntry("Esc", "cancel"),
		}
	} else {
		entries = []string{
			helpEntry("j/k", "move"),
			helpEntry("←→", "tabs"),
			helpEntry("Tab", "focus"),
			helpEntry("/", "search"),
			helpEntry("?", "content"),
			helpEntry("g", "group"),
			helpEntry("p", "pin"),
			helpEntry("e", "edit"),
			helpEntry("q", "quit"),
		}
	}

	return "  " + strings.Join(entries, "  ")
}

// injectBorderTitle places a styled title string into the top border line.
// Works by replacing a segment of the first line after the border corner.
func injectBorderTitle(box, title string) string {
	if title == "" {
		return box
	}
	lines := strings.Split(box, "\n")
	if len(lines) == 0 {
		return box
	}
	// Replace the top border line: ╭──...──╮ → ╭─ Title ─...──╮
	first := lines[0]
	cornerEnd := strings.Index(first, "─")
	if cornerEnd < 0 {
		return box
	}
	// Insert title after first border char
	insertAt := cornerEnd + len("─")
	lines[0] = first[:insertAt] + title + first[insertAt+len(title):]
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
	// Find a position near the right side to inject the footer
	cornerEnd := strings.Index(last, "─")
	if cornerEnd < 0 {
		return box
	}
	insertAt := cornerEnd + len("─")
	if insertAt+len(footer) < len(last) {
		lines[len(lines)-1] = last[:insertAt] + footer + last[insertAt+len(footer):]
	}
	return strings.Join(lines, "\n")
}
