package ui

import (
	"github.com/charmbracelet/lipgloss"
	"main/internal/ui/theme"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.PrimaryColor).
			MarginBottom(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(theme.ErrorColor)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(theme.SuccessColor)

	MutedStyle = lipgloss.NewStyle().
			Foreground(theme.MutedColor)

	InfoStyle = lipgloss.NewStyle().
			Foreground(theme.InfoColor)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.PrimaryColor).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
			Foreground(theme.MutedColor).
			Padding(0, 2)

	TimerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			Padding(1, 2).
			BorderForeground(theme.PrimaryColor)

	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(theme.PrimaryColor).
				Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(theme.TextColor).
			Background(theme.Surface0Color).
			Padding(0, 1)

	StatusBarErrorStyle = lipgloss.NewStyle().
				Foreground(theme.BaseColor).
				Background(theme.ErrorColor).
				Padding(0, 1)

	StatusBarSuccessStyle = lipgloss.NewStyle().
				Foreground(theme.BaseColor).
				Background(theme.SuccessColor).
				Padding(0, 1)
)
