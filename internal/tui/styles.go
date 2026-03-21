package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("236")).
		Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Padding(0, 2)

	listItemStyle = lipgloss.NewStyle().
		Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1)

	detailStyle = lipgloss.NewStyle().
		Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	borderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238"))

	statusIcon    = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
	statusIconOff = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗")
)
