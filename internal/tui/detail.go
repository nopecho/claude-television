package tui

import (
	"fmt"
	"strings"
)

func (m model) renderDetailContent(height int) string {
	ch := m.selectedChannel()
	if ch == nil || ch.Data == nil {
		return detailStyle.Render("No channel selected")
	}

	var content string
	switch m.detailTab {
	case TabSettings:
		content = m.renderSettingsTab(ch)
	case TabClaudeMD:
		content = m.renderClaudeMDTab(ch)
	case TabHooks:
		content = m.renderHooksTab(ch)
	case TabMCP:
		content = m.renderMCPTab(ch)
	case TabPlugins:
		content = m.renderPluginsTab(ch)
	case TabGit:
		content = m.renderGitTab(ch)
	case TabMemory:
		content = m.renderMemoryTab(ch)
	}

	// Apply scroll offset
	lines := strings.Split(content, "\n")
	if m.detailScroll >= len(lines) {
		m.detailScroll = len(lines) - 1
	}
	if m.detailScroll < 0 {
		m.detailScroll = 0
	}

	end := m.detailScroll + height
	if end > len(lines) {
		end = len(lines)
	}

	visible := lines[m.detailScroll:end]
	result := strings.Join(visible, "\n")

	// Show scroll indicator if there's more content
	if end < len(lines) {
		result += "\n" + labelStyle.Render(fmt.Sprintf("  ↓ %d more lines", len(lines)-end))
	}
	if m.detailScroll > 0 {
		result = labelStyle.Render(fmt.Sprintf("  ↑ %d lines above", m.detailScroll)) + "\n" + result
	}

	return result
}
