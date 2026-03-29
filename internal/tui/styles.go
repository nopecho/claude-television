package tui

import "github.com/charmbracelet/lipgloss"

// Color palette — Catppuccin Mocha inspired, enriched
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
	cyanColor    = lipgloss.AdaptiveColor{Light: "80", Dark: "#89DCEB"}
	baseColor    = lipgloss.AdaptiveColor{Light: "234", Dark: "#1E1E2E"}
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

// Header bar
var (
	headerAppStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	headerSepStyle = lipgloss.NewStyle().
			Foreground(overlayColor)

	headerChannelStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Bold(true)

	headerBadgeHealthy = lipgloss.NewStyle().
				Foreground(greenColor)

	headerBadgeWarning = lipgloss.NewStyle().
				Foreground(warningColor)

	headerBadgeError = lipgloss.NewStyle().
				Foreground(errorColor)

	headerCountStyle = lipgloss.NewStyle().
				Foreground(subtextColor)
)

// Tab bar
var (
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(baseColor).
			Background(accentColor).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(subtextColor).
				Padding(0, 1)

	tabBadgeStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)
)

// Content card styles
var (
	cardBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(overlayColor).
			Padding(0, 1)

	cardTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

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
				Background(surfaceColor).
				Padding(0, 1)

	channelSelectedPathStyle = lipgloss.NewStyle().
					Foreground(subtleColor).
					Background(surfaceColor).
					Padding(0, 1)

	groupHeaderStyle = lipgloss.NewStyle().
				Foreground(subtleColor).
				Padding(0, 1)

	groupDividerStyle = lipgloss.NewStyle().
				Foreground(overlayColor)

	channelBadgeStyle = lipgloss.NewStyle().
				Foreground(warningColor)
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
			Bold(true).
			Foreground(baseColor).
			Background(overlayColor).
			Padding(0, 1)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(subtextColor)

	helpScrollStyle = lipgloss.NewStyle().
			Foreground(subtleColor)
)

// Search
var searchStyle = lipgloss.NewStyle().
	Foreground(accentColor).
	Bold(true)

// Content search highlights
var (
	contentMatchStyle = lipgloss.NewStyle().
				Background(lipgloss.AdaptiveColor{Light: "222", Dark: "#F9E2AF"}).
				Foreground(lipgloss.Color("#1E1E2E"))

	contentMatchCurrentStyle = lipgloss.NewStyle().
					Background(lipgloss.AdaptiveColor{Light: "166", Dark: "#FAB387"}).
					Foreground(lipgloss.Color("#1E1E2E")).
					Bold(true)
)

// Detail content padding
var detailStyle = lipgloss.NewStyle().
	Padding(0, 1)

// Empty state
var (
	emptyArtStyle = lipgloss.NewStyle().
			Foreground(overlayColor)

	emptyMsgStyle = lipgloss.NewStyle().
			Foreground(subtextColor)

	emptyHintStyle = lipgloss.NewStyle().
			Foreground(subtleColor)
)

// Scrollbar
var (
	scrollTrackStyle = lipgloss.NewStyle().
				Foreground(overlayColor)

	scrollThumbStyle = lipgloss.NewStyle().
				Foreground(accentColor)
)
