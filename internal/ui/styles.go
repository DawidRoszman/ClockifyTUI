package ui

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryColor = lipgloss.Color("#7C3AED")
	SuccessColor = lipgloss.Color("#10B981")
	ErrorColor   = lipgloss.Color("#EF4444")
	MutedColor   = lipgloss.Color("#6B7280")
	InfoColor    = lipgloss.Color("#3B82F6")
	WarningColor = lipgloss.Color("#F59E0B")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Padding(0, 2)

	TimerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			Padding(1, 2).
			BorderForeground(PrimaryColor)

	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(PrimaryColor).
				Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(MutedColor).
			Padding(0, 1)

	StatusBarErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(ErrorColor).
				Padding(0, 1)

	StatusBarSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(SuccessColor).
				Padding(0, 1)
)
