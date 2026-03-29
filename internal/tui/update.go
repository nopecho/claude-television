package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/config"
)

// ansiEscRe strips ANSI escape sequences for plain-text searching.
var ansiEscRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching || m.contentSearching || m.grouping {
			return m.updateSearch(msg)
		}
		return m.updateNormal(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.syncViewportSize()
		m.syncDetailContent()
		return m, nil
	case tea.MouseMsg:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	case editorFinishedMsg:
		return m, nil
	}
	return m, nil
}

func (m *model) syncViewportSize() {
	listWidth := m.listWidth()
	detailWidth := m.width - listWidth - 4
	contentHeight := m.height - 5 // header + help + borders
	if contentHeight < 1 {
		contentHeight = 1
	}
	m.viewport.Width = detailWidth - 2
	if m.viewport.Width < 1 {
		m.viewport.Width = 1
	}
	tabBarHeight := 1
	m.viewport.Height = contentHeight - tabBarHeight
	if m.viewport.Height < 1 {
		m.viewport.Height = 1
	}
}

func (m model) listWidth() int {
	w := m.width * 25 / 100
	if w < 28 {
		w = 28
	}
	if w > 36 {
		w = 36
	}
	return w
}

func (m model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action := parseKey(msg)

	switch action {
	case keyQuit:
		return m, tea.Quit

	case keyTab, keyShiftTab:
		if m.focus == listPanel {
			m.focus = detailPanel
		} else {
			m.focus = listPanel
		}
		return m, nil

	case keyUp:
		if m.focus == listPanel {
			if m.channelCursor > 0 {
				m.channelCursor--
				m.syncDetailContent()
			}
		} else {
			m.viewport.LineUp(1)
		}

	case keyDown:
		if m.focus == listPanel {
			if m.channelCursor < len(m.filtered)-1 {
				m.channelCursor++
				m.syncDetailContent()
			}
		} else {
			m.viewport.LineDown(1)
		}

	case keyLeft:
		if m.focus == detailPanel {
			m.detailTab = (m.detailTab - 1 + DetailTab(len(detailTabNames))) % DetailTab(len(detailTabNames))
			m.syncDetailContent()
		}

	case keyRight:
		if m.focus == listPanel {
			m.focus = detailPanel
		} else {
			m.detailTab = (m.detailTab + 1) % DetailTab(len(detailTabNames))
			m.syncDetailContent()
		}

	case keyTab1, keyTab2, keyTab3, keyTab4, keyTab5, keyTab6, keyTab7, keyTab8:
		tabIdx := DetailTab(action - keyTab1)
		if tabIdx != m.detailTab {
			m.detailTab = tabIdx
			m.syncDetailContent()
		}
		return m, nil

	case keySlash:
		m.prevCursor = m.channelCursor
		m.searching = true
		m.searchInput.Placeholder = "Search channels..."
		m.searchInput.SetValue("")
		m.searchInput.Focus()
		return m, textinput.Blink

	case keyContentSearch:
		m.prevCursor = m.channelCursor
		m.contentSearching = true
		m.searchInput.Placeholder = "Search content..."
		m.searchInput.SetValue("")
		m.searchInput.Focus()
		return m, textinput.Blink

	case keyGroup:
		ch := m.selectedChannel()
		if ch != nil && !ch.IsGlobal {
			m.grouping = true
			m.searchInput.Placeholder = "Set group name..."
			m.searchInput.SetValue(ch.Group)
			m.searchInput.Focus()
			return m, textinput.Blink
		}

	case keyCmdEnter:
		ch := m.selectedChannel()
		if ch != nil {
			m.navigateTo = ch.Path
			return m, tea.Quit
		}

	case keyPin:
		ch := m.selectedChannel()
		if ch != nil {
			ch.Pinned = !ch.Pinned
			m.updatePins()
		}

	case keyEdit:
		ch := m.selectedChannel()
		if ch != nil {
			editor := m.cfg.Editor
			if editor == "" {
				editor = os.Getenv("EDITOR")
			}
			if editor != "" {
				target := editTargetForTab(ch, m.detailTab)
				c := exec.Command(editor, target)
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return editorFinishedMsg{err}
				})
			}
		}

	case keyScrollDown:
		m.viewport.HalfViewDown()

	case keyScrollUp:
		m.viewport.HalfViewUp()
	}

	return m, nil
}

func (m *model) syncDetailContent() {
	if m.contentSearching && m.searchInput.Value() != "" {
		m.applyContentSearch()
		return
	}
	content := m.renderDetailContentString()
	m.viewport.SetContent(content)
	m.viewport.GotoTop()
}

func (m *model) updatePins() {
	var pins []string
	for _, ch := range m.channels {
		if ch.Pinned {
			pins = append(pins, ch.Name)
		}
	}
	m.cfg.Channels.Pins = pins
	config.Save(m.cfg)
}

func (m model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searching = false
		m.contentSearching = false
		m.contentMatches = nil
		m.contentMatchIdx = 0
		m.grouping = false
		m.searchInput.Blur()
		m.searchInput.SetValue("")
		m.resetFilter()
		m.channelCursor = m.prevCursor
		if m.channelCursor >= len(m.filtered) {
			m.channelCursor = 0
		}
		m.syncDetailContent()
		return m, nil
	case "enter":
		if m.grouping {
			ch := m.selectedChannel()
			if ch != nil {
				ch.Group = m.searchInput.Value()
				m.updateGroups()
				m.sortChannels()
				m.resetFilter()
			}
		}
		m.searching = false
		m.contentSearching = false
		m.contentMatches = nil
		m.contentMatchIdx = 0
		m.grouping = false
		m.searchInput.Blur()
		m.syncDetailContent()
		return m, nil
	case "n":
		if m.contentSearching && len(m.contentMatches) > 0 {
			m.contentMatchIdx = (m.contentMatchIdx + 1) % len(m.contentMatches)
			m.renderContentWithHighlight()
			m.viewport.SetYOffset(m.contentMatches[m.contentMatchIdx])
			return m, nil
		}
	case "N":
		if m.contentSearching && len(m.contentMatches) > 0 {
			m.contentMatchIdx = (m.contentMatchIdx - 1 + len(m.contentMatches)) % len(m.contentMatches)
			m.renderContentWithHighlight()
			m.viewport.SetYOffset(m.contentMatches[m.contentMatchIdx])
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	if m.contentSearching {
		m.applyContentSearch()
	} else if !m.grouping {
		m.applySearch()
	}
	return m, cmd
}

func (m *model) applySearch() {
	query := m.searchInput.Value()
	if query == "" {
		m.resetFilter()
		return
	}

	results := channel.FuzzySearch(m.channels, query)
	if len(results) == 0 {
		m.filtered = []int{}
		m.channelCursor = 0
		return
	}

	indexMap := make(map[string]int, len(m.channels))
	for i, ch := range m.channels {
		indexMap[ch.ID] = i
	}

	m.filtered = make([]int, 0, len(results))
	for _, r := range results {
		if idx, ok := indexMap[r.ID]; ok {
			m.filtered = append(m.filtered, idx)
		}
	}
	m.channelCursor = 0
}

func (m *model) applyContentSearch() {
	query := strings.ToLower(m.searchInput.Value())
	content := m.renderDetailContentString()
	lines := strings.Split(content, "\n")

	m.contentMatches = nil
	if query == "" {
		m.viewport.SetContent(content)
		return
	}

	for i, line := range lines {
		plain := strings.ToLower(ansiEscRe.ReplaceAllString(line, ""))
		if strings.Contains(plain, query) {
			m.contentMatches = append(m.contentMatches, i)
		}
	}

	m.contentMatchIdx = 0
	m.highlightAndSetContent(lines)
	if len(m.contentMatches) > 0 {
		m.viewport.SetYOffset(m.contentMatches[0])
	}
}

func (m *model) renderContentWithHighlight() {
	if len(m.contentMatches) == 0 {
		m.viewport.SetContent(m.renderDetailContentString())
		return
	}
	lines := strings.Split(m.renderDetailContentString(), "\n")
	m.highlightAndSetContent(lines)
}

func (m *model) highlightAndSetContent(lines []string) {
	currentLine := -1
	if m.contentMatchIdx < len(m.contentMatches) {
		currentLine = m.contentMatches[m.contentMatchIdx]
	}
	for _, lineIdx := range m.contentMatches {
		if lineIdx >= len(lines) {
			continue
		}
		plain := ansiEscRe.ReplaceAllString(lines[lineIdx], "")
		if lineIdx == currentLine {
			lines[lineIdx] = contentMatchCurrentStyle.Render(plain)
		} else {
			lines[lineIdx] = contentMatchStyle.Render(plain)
		}
	}
	m.viewport.SetContent(strings.Join(lines, "\n"))
}

func (m *model) updateGroups() {
	groups := map[string][]string{}
	for _, ch := range m.channels {
		if ch.Group != "" && !ch.IsGlobal {
			groups[ch.Group] = append(groups[ch.Group], ch.Name)
		}
	}
	m.cfg.Channels.Groups = groups
	config.Save(m.cfg)
}

func editTargetForTab(ch *channel.Channel, tab DetailTab) string {
	switch tab {
	case TabClaudeMD:
		return filepath.Join(ch.Path, "CLAUDE.md")
	default:
		return filepath.Join(ch.Path, ".claude", "settings.json")
	}
}
