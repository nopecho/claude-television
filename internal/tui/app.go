package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nopecho/claude-television/internal/claude"
	"github.com/nopecho/claude-television/internal/scanner"
)

type DashboardData struct {
	Settings      *claude.Settings
	LocalSettings *claude.Settings
	ClaudeMD      *claude.ClaudeMD
	Plugins       []claude.Plugin
	LocalSkills   []claude.Skill
	Hooks         []claude.HookDetail
	ProjectsMeta  []claude.ProjectMeta
	Projects      []scanner.Project
}

type skillItem struct {
	display   string
	isPlugin  bool
	pluginIdx int
	skillIdx  int
}

type model struct {
	data       DashboardData
	activeTab  TabID
	cursor     int
	width      int
	height     int
	skillCache []skillItem
}

func NewModel(data DashboardData) model {
	m := model{
		data:      data,
		activeTab: TabGlobal,
	}
	m.skillCache = m.buildSkillItems()
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab", "right", "l":
			m.activeTab = (m.activeTab + 1) % TabID(len(tabNames))
			m.cursor = 0
		case "shift+tab", "left", "h":
			m.activeTab = (m.activeTab - 1 + TabID(len(tabNames))) % TabID(len(tabNames))
			m.cursor = 0
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			max := m.listLen() - 1
			if m.cursor < max {
				m.cursor++
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	header := titleStyle.Render("📺 claude-television") + "\n\n"
	tabs := renderTabs(m.activeTab) + "\n\n"

	listWidth := m.width*35/100 - 4
	detailWidth := m.width*65/100 - 4
	contentHeight := m.height - 8

	list := borderStyle.Width(listWidth).Height(contentHeight).Render(m.renderList())
	detail := borderStyle.Width(detailWidth).Height(contentHeight).Render(m.renderDetail())

	content := lipgloss.JoinHorizontal(lipgloss.Top, list, detail)
	help := helpStyle.Render("\n  ↑↓/jk navigate  ←→/Tab switch  q quit")

	return header + tabs + content + help
}

func (m model) listLen() int {
	switch m.activeTab {
	case TabGlobal:
		return 3
	case TabProjects:
		return len(m.data.Projects)
	case TabSkills:
		return len(m.skillCache)
	case TabHooks:
		return len(m.data.Hooks)
	}
	return 0
}

func (m model) renderList() string {
	switch m.activeTab {
	case TabGlobal:
		return m.renderGlobalList()
	case TabProjects:
		return m.renderProjectsList()
	case TabSkills:
		return m.renderSkillsList()
	case TabHooks:
		return m.renderHooksList()
	}
	return ""
}

func (m model) renderDetail() string {
	switch m.activeTab {
	case TabGlobal:
		return m.renderGlobalDetail()
	case TabProjects:
		return m.renderProjectsDetail()
	case TabSkills:
		return m.renderSkillsDetail()
	case TabHooks:
		return m.renderHooksDetail()
	}
	return ""
}

func Run(data DashboardData) error {
	p := tea.NewProgram(NewModel(data), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
