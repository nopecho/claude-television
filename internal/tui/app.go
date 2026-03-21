package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/config"
)

type DetailTab int

const (
	TabSettings DetailTab = iota
	TabClaudeMD
	TabHooks
	TabMCP
	TabPlugins
	TabGit
	TabMemory
)

var detailTabNames = []string{"Settings", "CLAUDE.md", "Hooks", "MCP", "Plugins", "Git", "Memory"}

type model struct {
	channels      []channel.Channel
	cfg           *config.Config
	channelCursor int
	detailTab     DetailTab
	detailScroll  int
	width         int
	height        int
	searching     bool
	searchQuery   string
	filtered      []int
	navigateTo    string
}

func newModel(channels []channel.Channel, cfg *config.Config) model {
	m := model{
		channels: channels,
		cfg:      cfg,
	}
	m.sortChannels()
	m.resetFilter()
	return m
}

func (m *model) sortChannels() {
	pinned := make([]channel.Channel, 0)
	grouped := make(map[string][]channel.Channel)
	ungrouped := make([]channel.Channel, 0)
	var groupOrder []string

	for _, ch := range m.channels {
		if ch.Pinned {
			pinned = append(pinned, ch)
		} else if ch.Group != "" {
			if _, exists := grouped[ch.Group]; !exists {
				groupOrder = append(groupOrder, ch.Group)
			}
			grouped[ch.Group] = append(grouped[ch.Group], ch)
		} else {
			ungrouped = append(ungrouped, ch)
		}
	}

	sorted := make([]channel.Channel, 0, len(m.channels))
	sorted = append(sorted, pinned...)
	for _, g := range groupOrder {
		sorted = append(sorted, grouped[g]...)
	}
	sorted = append(sorted, ungrouped...)
	m.channels = sorted
}

func (m *model) resetFilter() {
	m.filtered = make([]int, len(m.channels))
	for i := range m.channels {
		m.filtered[i] = i
	}
}

func (m model) selectedChannel() *channel.Channel {
	if len(m.filtered) == 0 || m.channelCursor >= len(m.filtered) {
		return nil
	}
	idx := m.filtered[m.channelCursor]
	return &m.channels[idx]
}

func (m model) Init() tea.Cmd {
	return nil
}

type editorFinishedMsg struct{ err error }

func RunChannels(channels []channel.Channel, cfg *config.Config) error {
	m := newModel(channels, cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return err
	}
	if final, ok := result.(model); ok && final.navigateTo != "" {
		fmt.Print(final.navigateTo)
	}
	return nil
}
