package tui

import "github.com/charmbracelet/lipgloss"

// Color palette.
var (
	primaryColor   = lipgloss.Color("39")  // Azure blue
	secondaryColor = lipgloss.Color("208") // Orange
	successColor   = lipgloss.Color("82")  // Green
	errorColor     = lipgloss.Color("196") // Red
	mutedColor     = lipgloss.Color("241") // Gray
	highlightColor = lipgloss.Color("212") // Pink
)

// Styles for the TUI.
var (
	// Title style for headers.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Subtitle style.
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginBottom(1)

	// Selected item style.
	SelectedStyle = lipgloss.NewStyle().
			Foreground(highlightColor).
			Bold(true)

	// Current item indicator style.
	CurrentStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// Normal item style.
	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	// Muted style for secondary text.
	MutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Error style.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Success style.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor)

	// Help style.
	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Box style for sections.
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1)

	// Header box style.
	HeaderBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1).
			MarginBottom(1)

	// Cursor style.
	CursorStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	// Tab style.
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Underline(true)

	// Inactive tab style.
	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	// Spinner style.
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	// Status bar style.
	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1)
)
