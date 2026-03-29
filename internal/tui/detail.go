package tui

// renderDetailContentString produces the raw content string for the viewport.
func (m model) renderDetailContentString() string {
	ch := m.selectedChannel()
	if ch == nil || ch.Data == nil {
		return detailStyle.Render("No channel selected")
	}

	switch m.detailTab {
	case TabSettings:
		return m.renderSettingsTab(ch)
	case TabClaudeMD:
		return m.renderClaudeMDTab(ch)
	case TabHooks:
		return m.renderHooksTab(ch)
	case TabMCP:
		return m.renderMCPTab(ch)
	case TabPlugins:
		return m.renderPluginsTab(ch)
	case TabHealth:
		return m.renderHealthTab(ch)
	case TabGit:
		return m.renderGitTab(ch)
	case TabMemory:
		return m.renderMemoryTab(ch)
	}
	return ""
}

// detailContentWidth returns the available width for detail content cards.
func (m model) detailContentWidth() int {
	listWidth := m.listWidth()
	w := m.width - listWidth - 6 // borders + padding
	if w < 20 {
		w = 20
	}
	return w
}
