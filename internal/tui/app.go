package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nopecho/claude-television/internal/channel"
	"github.com/nopecho/claude-television/internal/config"
)

type DetailTab int

const (
	TabSettings DetailTab = iota
	TabClaudeMD
	TabHooks
	TabMCP
	TabGit
	TabMemory
)

var detailTabNames = []string{"Settings", "CLAUDE.md", "Hooks", "MCP", "Git", "Memory"}

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
		}
	case keyEdit:
		ch := m.selectedChannel()
		if ch != nil {
			m.openEditor(ch)
		}
	}
	return m, nil
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
	m.filtered = make([]int, 0)
	for _, r := range results {
		for i, ch := range m.channels {
			if ch.ID == r.ID {
				m.filtered = append(m.filtered, i)
				break
			}
		}
	}
	m.channelCursor = 0
}

func (m model) openEditor(ch *channel.Channel) {
	editor := m.cfg.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		return
	}
	var target string
	switch m.detailTab {
	case TabSettings:
		target = ch.Path + "/.claude/settings.json"
	case TabClaudeMD:
		target = ch.Path + "/CLAUDE.md"
	default:
		target = ch.Path + "/.claude/settings.json"
	}
	cmd := exec.Command(editor, target)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

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

	help := helpStyle.Render("  j/k move  ←→/Tab switch tab  / search  Alt+Enter cd  p pin  e edit  q quit")

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
