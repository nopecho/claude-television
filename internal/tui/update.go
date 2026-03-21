package tui

import (
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/config"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching {
			return m.updateSearch(msg)
		}
		return m.updateNormal(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case editorFinishedMsg:
		// Editor finished, no action needed - TUI resumes automatically
		return m, nil
	}
	return m, nil
}

func (m model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action := parseKey(msg)
	switch action {
	case keyQuit:
		return m, tea.Quit
	case keyUp:
		if m.channelCursor > 0 {
			m.channelCursor--
			m.detailScroll = 0
		}
	case keyDown:
		if m.channelCursor < len(m.filtered)-1 {
			m.channelCursor++
			m.detailScroll = 0
		}
	case keyTab, keyRight:
		m.detailTab = (m.detailTab + 1) % DetailTab(len(detailTabNames))
		m.detailScroll = 0
	case keyShiftTab, keyLeft:
		m.detailTab = (m.detailTab - 1 + DetailTab(len(detailTabNames))) % DetailTab(len(detailTabNames))
		m.detailScroll = 0
	case keySlash:
		m.searching = true
		m.searchQuery = ""
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
		m.detailScroll += 10
	case keyScrollUp:
		m.detailScroll -= 10
		if m.detailScroll < 0 {
			m.detailScroll = 0
		}
	}
	return m, nil
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
		m.searchQuery = ""
		m.resetFilter()
		m.channelCursor = 0
	case "enter":
		m.searching = false
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.applySearch()
		}
	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
			m.applySearch()
		}
	}
	return m, nil
}

func (m *model) applySearch() {
	if m.searchQuery == "" {
		m.resetFilter()
		return
	}
	results := channel.FuzzySearch(m.channels, m.searchQuery)
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
	m.detailScroll = 0
}

func editTargetForTab(ch *channel.Channel, tab DetailTab) string {
	switch tab {
	case TabClaudeMD:
		return filepath.Join(ch.Path, "CLAUDE.md")
	default:
		return filepath.Join(ch.Path, ".claude", "settings.json")
	}
}
