package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"main/internal/ui/theme"
)

type StatusBarComponent struct {
	message string
	msgType StatusType
	width   int
}

type StatusType int

const (
	StatusNormal StatusType = iota
	StatusSuccess
	StatusError
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(theme.TextColor).
			Background(theme.Surface0Color).
			Padding(0, 1)

	statusBarErrorStyle = lipgloss.NewStyle().
				Foreground(theme.BaseColor).
				Background(theme.RedColor).
				Padding(0, 1)

	statusBarSuccessStyle = lipgloss.NewStyle().
				Foreground(theme.BaseColor).
				Background(theme.GreenColor).
				Padding(0, 1)
)

func NewStatusBar() *StatusBarComponent {
	return &StatusBarComponent{
		message: "Ready",
		msgType: StatusNormal,
	}
}

func (c *StatusBarComponent) SetMessage(msg string, msgType StatusType) {
	c.message = msg
	c.msgType = msgType
}

func (c *StatusBarComponent) SetWidth(width int) {
	c.width = width
}

func (c *StatusBarComponent) View() string {
	style := statusBarStyle
	switch c.msgType {
	case StatusSuccess:
		style = statusBarSuccessStyle
	case StatusError:
		style = statusBarErrorStyle
	}

	if c.width > 0 {
		style = style.Width(c.width)
	}

	return style.Render(c.message)
}

func (c *StatusBarComponent) SetError(err error) {
	c.SetMessage(fmt.Sprintf("Error: %v", err), StatusError)
}

func (c *StatusBarComponent) SetSuccess(msg string) {
	c.SetMessage(msg, StatusSuccess)
}

func (c *StatusBarComponent) SetInfo(msg string) {
	c.SetMessage(msg, StatusNormal)
}
