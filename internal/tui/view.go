package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/nopecho/claude-television/internal/channel"
)

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	listWidth := m.width * 22 / 100
	if listWidth < 20 {
		listWidth = 20
	}
	detailWidth := m.width - listWidth - 4
	contentHeight := m.height - 5

	header := titleStyle.Render("ctv") + " "
	if m.searching {
		header += searchStyle.Render("/ " + m.searchQuery + "█")
	} else if m.contentSearching {
		header += searchStyle.Render("? " + m.searchQuery + "█") + " " + labelStyle.Render("(content)")
	} else if m.grouping {
		header += searchStyle.Render("group: " + m.searchQuery + "█") + " " + labelStyle.Render("(enter to set, empty to clear)")
	} else {
		header += m.renderSummary()
	}
	header += "\n"

	listContent := m.renderChannelList(contentHeight)
	list := borderStyle.Width(listWidth).Height(contentHeight).Render(listContent)

	tabBar := m.renderDetailTabs()
	detailContent := m.renderDetailContent(contentHeight - 2)
	detail := borderStyle.Width(detailWidth).Height(contentHeight).Render(tabBar + "\n" + detailContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, list, detail)

	help := helpStyle.Render("  j/k move  ←→/Tab switch  / search  ? content  g group  p pin  e edit  Alt+Enter cd  q quit")

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
	var tabs string
	for i, name := range detailTabNames {
		if DetailTab(i) == m.detailTab {
			tabs += activeTabStyle.Render("[" + name + "]")
		} else {
			tabs += inactiveTabStyle.Render(" " + name + " ")
		}
	}
	return tabs
}
