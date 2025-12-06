package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"main/internal/api"
	"main/internal/domain"
	"main/internal/ui/components"
	"main/internal/ui/theme"
)

type TimerView struct {
	timerComponent    *components.TimerComponent
	projectSelector   *components.ProjectSelectorComponent
	showSelector      bool
	width             int
	height            int
}

func NewTimerView(timerState *domain.TimerState) *TimerView {
	return &TimerView{
		timerComponent:  components.NewTimerComponent(timerState),
		projectSelector: components.NewProjectSelector(),
		showSelector:    false,
	}
}

func (v *TimerView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.projectSelector.SetSize(width, height)
}

func (v *TimerView) SetProjects(projects []api.Project) {
	v.projectSelector.SetProjects(projects)
}

func (v *TimerView) SetProjectMap(projects map[string]string) {
	v.timerComponent.SetProjectMap(projects)
}

func (v *TimerView) SetTagMap(tags map[string]string) {
	v.timerComponent.SetTagMap(tags)
}

func (v *TimerView) ShowProjectSelector() {
	v.showSelector = true
}

func (v *TimerView) HideProjectSelector() {
	v.showSelector = false
}

func (v *TimerView) IsShowingSelector() bool {
	return v.showSelector
}

func (v *TimerView) GetProjectSelector() *components.ProjectSelectorComponent {
	return v.projectSelector
}

func (v *TimerView) GetTimerComponent() *components.TimerComponent {
	return v.timerComponent
}

func (v *TimerView) Update(msg tea.Msg) (*TimerView, tea.Cmd) {
	return v, nil
}

func (v *TimerView) View() string {
	if v.showSelector {
		return v.projectSelector.View()
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.LavenderColor).
		MarginBottom(1)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.Subtext0Color)

	content := titleStyle.Render("⏱  Timer") + "\n\n"
	content += v.timerComponent.View() + "\n\n"

	if v.timerComponent.View() != "" {
		helpText := "s: start timer"
		if v.timerComponent.IsRunning() {
			helpText += " | x: stop timer | d: edit description"
		}
		helpText += " | p: select project"
		content += mutedStyle.Render(helpText)
	}

	return content
}
