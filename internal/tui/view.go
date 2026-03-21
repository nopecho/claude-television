package tui

import (
	"fmt"
	"strings"

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
	}
	header += "\n"

	listContent := m.renderChannelList(contentHeight)
	list := borderStyle.Width(listWidth).Height(contentHeight).Render(listContent)

	tabBar := m.renderDetailTabs()
	detailContent := m.renderDetailContent(contentHeight - 2)
	detail := borderStyle.Width(detailWidth).Height(contentHeight).Render(tabBar + "\n" + detailContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, list, detail)

	help := helpStyle.Render("  j/k move  ←→/Tab switch tab  / search  Ctrl+d/u scroll  Alt+Enter cd  p pin  e edit  q quit")

	return header + content + "\n" + help
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
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func kv(key, value string) string {
	return fmt.Sprintf("  %s  %s", labelStyle.Render(key+":"), valueStyle.Render(value))
}

func section(title string) string {
	return "\n" + headerStyle.Render("  "+title) + "\n"
}

func indent(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = "    " + l
	}
	return strings.Join(lines, "\n")
}
