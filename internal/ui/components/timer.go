package components

import (
	"github.com/charmbracelet/lipgloss"
	"main/internal/domain"
	"main/internal/ui/theme"
)

type TimerComponent struct {
	timerState         *domain.TimerState
	projects           map[string]string
	tags               map[string]string
	editingDescription bool
	editBuffer         string
}

var (
	timerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			Padding(1, 2).
			BorderForeground(theme.PrimaryColor)

	timerRunningColor = theme.GreenColor
	timerStoppedColor = theme.Overlay0Color
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

func (c *TimerComponent) SetTagMap(tags map[string]string) {
	c.tags = tags
}

func (c *TimerComponent) IsRunning() bool {
	return c.timerState.IsRunning
}

func (c *TimerComponent) StartEditingDescription() {
	if c.timerState.IsRunning {
		c.editingDescription = true
		c.editBuffer = c.timerState.Description
	}
}

func (c *TimerComponent) IsEditingDescription() bool {
	return c.editingDescription
}

func (c *TimerComponent) CancelEditingDescription() {
	c.editingDescription = false
	c.editBuffer = ""
}

func (c *TimerComponent) ClearEditState() {
	c.editingDescription = false
	c.editBuffer = ""
}

func (c *TimerComponent) GetEditedDescription() string {
	return c.editBuffer
}

func (c *TimerComponent) AddCharToEdit(char rune) {
	if c.editingDescription && len(c.editBuffer) < 255 {
		c.editBuffer += string(char)
	}
}

func (c *TimerComponent) DeleteCharFromEdit() {
	if c.editingDescription && len(c.editBuffer) > 0 {
		c.editBuffer = c.editBuffer[:len(c.editBuffer)-1]
	}
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

	if c.editingDescription {
		content += lipgloss.NewStyle().Bold(true).Render("Description: ") + "\n"

		inputStyle := lipgloss.NewStyle().
			Foreground(theme.GreenColor).
			Bold(true)

		content += "  " + inputStyle.Render(c.editBuffer) + "█\n"
		content += "\n" + lipgloss.NewStyle().
			Foreground(theme.Subtext0Color).
			Render("  enter: save | esc: cancel") + "\n"
	} else {
		if c.timerState.Description != "" {
			content += lipgloss.NewStyle().Bold(true).Render("Description: ") + c.timerState.Description + "\n"
		}
	}

	if c.timerState.ProjectID != nil {
		projectName := *c.timerState.ProjectID
		if name, ok := c.projects[*c.timerState.ProjectID]; ok {
			projectName = name
		}
		content += lipgloss.NewStyle().Bold(true).Render("Project: ") + projectName + "\n"
	}

	if len(c.timerState.TagIDs) > 0 {
		tagNames := []string{}
		for _, tagID := range c.timerState.TagIDs {
			if name, ok := c.tags[tagID]; ok {
				tagNames = append(tagNames, name)
			} else {
				tagNames = append(tagNames, tagID)
			}
		}

		if len(tagNames) > 0 {
			tagStyle := lipgloss.NewStyle().
				Foreground(theme.MauveColor).
				Italic(true)

			tagsStr := ""
			for i, name := range tagNames {
				if i > 0 {
					tagsStr += ", "
				}
				tagsStr += name
			}

			content += lipgloss.NewStyle().Bold(true).Render("Tags: ") +
				tagStyle.Render(tagsStr) + "\n"
		}
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
