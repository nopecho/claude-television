package tui

import "github.com/charmbracelet/lipgloss"

// Color palette — Catppuccin Mocha inspired
var (
	accentColor  = lipgloss.AdaptiveColor{Light: "97", Dark: "#7C6FAF"}
	greenColor   = lipgloss.AdaptiveColor{Light: "65", Dark: "#5B8A72"}
	subtleColor  = lipgloss.AdaptiveColor{Light: "243", Dark: "#6C7086"}
	surfaceColor = lipgloss.AdaptiveColor{Light: "236", Dark: "#313244"}
	overlayColor = lipgloss.AdaptiveColor{Light: "239", Dark: "#45475A"}
	textColor    = lipgloss.AdaptiveColor{Light: "189", Dark: "#CDD6F4"}
	subtextColor = lipgloss.AdaptiveColor{Light: "145", Dark: "#A6ADC8"}
	warningColor = lipgloss.AdaptiveColor{Light: "222", Dark: "#F9E2AF"}
	errorColor   = lipgloss.AdaptiveColor{Light: "211", Dark: "#F38BA8"}
)

// Panel borders
var (
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(accentColor)

	unfocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(overlayColor)
)

// Panel titles
var (
	focusedTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor)

	unfocusedTitleStyle = lipgloss.NewStyle().
				Foreground(subtleColor)
)

// Header / app title
var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(accentColor).
	Padding(0, 1)

// Tab bar
var (
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(subtextColor).
				Padding(0, 1)

	tabUnderlineStyle = lipgloss.NewStyle().
				Foreground(accentColor)
)

// Content
var (
	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor)

	sectionBarStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	labelStyle = lipgloss.NewStyle().
			Foreground(subtextColor)

	valueStyle = lipgloss.NewStyle().
			Foreground(textColor)
)

// Channel list
var (
	channelItemStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Padding(0, 1)

	channelSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor).
				Padding(0, 1)

	groupHeaderStyle = lipgloss.NewStyle().
				Foreground(subtleColor).
				Bold(true).
				Padding(0, 1)
)

// Status icons
var (
	statusHealthy = lipgloss.NewStyle().Foreground(greenColor).Render("●")
	statusWarning = lipgloss.NewStyle().Foreground(warningColor).Render("○")
	statusError   = lipgloss.NewStyle().Foreground(errorColor).Render("✕")
	pinIcon       = lipgloss.NewStyle().Foreground(warningColor).Render("★")
)

// Help bar
var (
	helpKeyStyle = lipgloss.NewStyle().
			Foreground(accentColor)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(subtextColor)
)

// Search
var searchStyle = lipgloss.NewStyle().
	Foreground(accentColor).
	Bold(true)

// Detail content padding
var detailStyle = lipgloss.NewStyle().
	Padding(0, 1)
