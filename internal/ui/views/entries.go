package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"main/internal/api"
	"main/internal/ui/components"
	"main/internal/ui/theme"
)

type EntriesView struct {
	entriesComponent *components.EntriesComponent
	width            int
	height           int
}

func NewEntriesView() *EntriesView {
	return &EntriesView{
		entriesComponent: components.NewEntriesComponent(),
	}
}

func (v *EntriesView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.entriesComponent.SetSize(width, height)
}

func (v *EntriesView) SetEntries(entries []api.TimeEntry) {
	v.entriesComponent.SetEntries(entries)
}

func (v *EntriesView) SetProjects(projects map[string]string) {
	v.entriesComponent.SetProjects(projects)
}

func (v *EntriesView) SetTasks(tasks map[string]string) {
	v.entriesComponent.SetTasks(tasks)
}

func (v *EntriesView) SetTags(tags map[string]string) {
	v.entriesComponent.SetTags(tags)
}

func (v *EntriesView) GetViewMode() components.EntriesViewMode {
	return v.entriesComponent.GetViewMode()
}

func (v *EntriesView) ToggleViewMode() {
	v.entriesComponent.ToggleViewMode()
}

func (v *EntriesView) MoveUp() {
	v.entriesComponent.PrevItem()
}

func (v *EntriesView) MoveDown() {
	v.entriesComponent.NextItem()
}

func (v *EntriesView) Update(msg tea.Msg) (*EntriesView, tea.Cmd) {
	return v, nil
}

func (v *EntriesView) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.LavenderColor).
		MarginBottom(1)

	viewModeStr := "Today"
	if v.entriesComponent.GetViewMode() == components.ViewThisWeek {
		viewModeStr = "This Week"
	}

	content := titleStyle.Render("📋 Time Entries - "+viewModeStr) + "\n\n"
	content += v.entriesComponent.View()

	return content
}
