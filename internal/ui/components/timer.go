package components

import (
	"github.com/charmbracelet/lipgloss"
	"main/internal/domain"
)

type TimerComponent struct {
	timerState *domain.TimerState
	projects   map[string]string
}

var (
	timerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			Padding(1, 2).
			BorderForeground(lipgloss.Color("#7C3AED"))

	timerRunningColor = lipgloss.Color("#10B981")
	timerStoppedColor = lipgloss.Color("#6B7280")
)

func NewTimerComponent(state *domain.TimerState) *TimerComponent {
	return &TimerComponent{
		timerState: state,
		projects:   make(map[string]string),
	}
}

func (c *TimerComponent) SetProjectMap(projects map[string]string) {
	c.projects = projects
}

func (c *TimerComponent) View() string {
	if c.timerState.IsRunning {
		return c.renderRunningTimer()
	}
	return c.renderStoppedTimer()
}

func (c *TimerComponent) renderRunningTimer() string {
	elapsed := c.timerState.GetElapsedDuration()
	durationStr := domain.FormatDuration(elapsed)

	statusStyle := lipgloss.NewStyle().
		Foreground(timerRunningColor).
		Bold(true)

	content := statusStyle.Render("⏱ RUNNING") + "\n\n"
	content += lipgloss.NewStyle().Bold(true).Render("Time: ") + durationStr + "\n"

	if c.timerState.Description != "" {
		content += lipgloss.NewStyle().Bold(true).Render("Description: ") + c.timerState.Description + "\n"
	}

	if c.timerState.ProjectID != nil {
		projectName := *c.timerState.ProjectID
		if name, ok := c.projects[*c.timerState.ProjectID]; ok {
			projectName = name
		}
		content += lipgloss.NewStyle().Bold(true).Render("Project: ") + projectName + "\n"
	}

	return timerBoxStyle.Render(content)
}

func (c *TimerComponent) renderStoppedTimer() string {
	statusStyle := lipgloss.NewStyle().
		Foreground(timerStoppedColor)

	content := statusStyle.Render("⏸ STOPPED") + "\n\n"
	content += lipgloss.NewStyle().Foreground(timerStoppedColor).Render("No timer running")

	return timerBoxStyle.Render(content)
}
