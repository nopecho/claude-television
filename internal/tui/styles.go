package tui

import "github.com/charmbracelet/lipgloss"

var (
	accentColor  = lipgloss.Color("205")
	dimColor     = lipgloss.Color("241")
	borderColor  = lipgloss.Color("238")
	healthyColor = lipgloss.Color("42")
	warningColor = lipgloss.Color("214")
	errorColor   = lipgloss.Color("196")
	groupColor   = lipgloss.Color("69")
	pinColor     = lipgloss.Color("220")
	tabBgColor   = lipgloss.Color("236")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor)

	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	channelItemStyle = lipgloss.NewStyle().
				Padding(0, 1)

	channelSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor).
				Padding(0, 1)

	groupHeaderStyle = lipgloss.NewStyle().
				Foreground(groupColor).
				Bold(true).
				Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Background(tabBgColor).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(dimColor).
				Padding(0, 1)

	detailStyle = lipgloss.NewStyle().
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	statusHealthy = lipgloss.NewStyle().Foreground(healthyColor).Render("●")
	statusWarning = lipgloss.NewStyle().Foreground(warningColor).Render("○")
	statusError   = lipgloss.NewStyle().Foreground(errorColor).Render("✕")
	pinIcon       = lipgloss.NewStyle().Foreground(pinColor).Render("★")

	searchStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)
)
