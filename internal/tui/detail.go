package tui

func (m model) renderDetailContent(height int) string {
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
	case TabGit:
		return m.renderGitTab(ch)
	case TabMemory:
		return m.renderMemoryTab(ch)
	}
	return ""
}
