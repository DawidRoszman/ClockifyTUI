package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	SwitchToTimer     key.Binding
	SwitchToEntries   key.Binding
	SwitchToReports   key.Binding
	StartTimer        key.Binding
	StopTimer         key.Binding
	SelectProject     key.Binding
	EditDescription   key.Binding
	Refresh           key.Binding
	Quit              key.Binding
	Up                key.Binding
	Down              key.Binding
	Left              key.Binding
	Right             key.Binding
	Enter             key.Binding
	Back              key.Binding
	Help              key.Binding
	ToggleView        key.Binding
	Space             key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		SwitchToTimer: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "timer"),
		),
		SwitchToEntries: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "entries"),
		),
		SwitchToReports: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "reports"),
		),
		StartTimer: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "start timer"),
		),
		StopTimer: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "stop timer"),
		),
		SelectProject: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "select project"),
		),
		EditDescription: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "edit description"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		ToggleView: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle selection"),
		),
	}
}
