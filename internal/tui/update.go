package tui

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/config"
)

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
	contentHeight := m.height - 7 // header + tab bar + help + borders
	if contentHeight < 1 {
		contentHeight = 1
	}
	m.viewport.Width = detailWidth - 2
	m.viewport.Height = contentHeight
}

func (m model) listWidth() int {
	w := m.width * 25 / 100
	if w < 24 {
		w = 24
	}
	if w > 30 {
		w = 30
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
		m.detailTab = (m.detailTab - 1 + DetailTab(len(detailTabNames))) % DetailTab(len(detailTabNames))
		m.syncDetailContent()

	case keyRight:
		m.detailTab = (m.detailTab + 1) % DetailTab(len(detailTabNames))
		m.syncDetailContent()

	case keySlash:
		m.searching = true
		m.searchInput.Placeholder = "Search channels..."
		m.searchInput.SetValue("")
		m.searchInput.Focus()
		return m, textinput.Blink

	case keyContentSearch:
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
		m.grouping = false
		m.searchInput.Blur()
		m.searchInput.SetValue("")
		m.resetFilter()
		m.channelCursor = 0
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
		m.grouping = false
		m.searchInput.Blur()
		m.syncDetailContent()
		return m, nil
	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		if !m.grouping {
			m.applySearch()
		}
		return m, cmd
	}
}

func (m *model) applySearch() {
	query := m.searchInput.Value()
	if query == "" {
		m.resetFilter()
		return
	}

	var results []channel.Channel
	if m.contentSearching {
		results = channel.ContentSearch(m.channels, query)
	} else {
		results = channel.FuzzySearch(m.channels, query)
	}

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
