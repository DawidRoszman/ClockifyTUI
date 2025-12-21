package ui

import (
	"fmt"
	"time"

	"main/internal/api"
	"main/internal/cache"
	"main/internal/domain"
	"main/internal/ui/components"
	"main/internal/ui/theme"
	"main/internal/ui/views"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	timerService   *domain.TimerService
	entryService   *domain.TimeEntryService
	reportService  *domain.ReportService
	projectService *domain.ProjectService
	tagService     *domain.TagService

	currentView ViewType
	width       int
	height      int

	timerView   *views.TimerView
	entriesView *views.EntriesView
	reportsView *views.ReportsView
	statusBar   *components.StatusBarComponent

	projects    []api.Project
	entries     []api.TimeEntry
	tags        []api.Tag
	projectsMap map[string]string
	tasksMap    map[string]string
	tagsMap     map[string]string

	showHelp  bool
	isLoading bool
	err       error

	keys KeyMap
}

func NewApp(client *api.Client) *App {
	cacheInstance := cache.NewCache(5 * time.Minute)
	timerState := domain.NewTimerState()
	timerService := domain.NewTimerService(client, timerState)

	return &App{
		timerService:   timerService,
		entryService:   domain.NewTimeEntryService(client),
		reportService:  domain.NewReportService(client),
		projectService: domain.NewProjectService(client, cacheInstance),
		tagService:     domain.NewTagService(client, cacheInstance),
		currentView:    TimerView,
		timerView:      views.NewTimerView(timerState),
		entriesView:    views.NewEntriesView(),
		reportsView:    views.NewReportsView(),
		statusBar:      components.NewStatusBar(),
		projectsMap:    make(map[string]string),
		tasksMap:       make(map[string]string),
		tagsMap:        make(map[string]string),
		keys:           DefaultKeyMap(),
	}
}

func (m App) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		m.loadCurrentTimer,
		m.loadProjects,
		m.loadTags,
	)
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusBar.SetWidth(m.width)
		m.timerView.SetSize(m.width, m.height)
		m.entriesView.SetSize(m.width, m.height)
		m.reportsView.SetSize(m.width, m.height)
		return m, nil

	case tea.KeyMsg:
		if m.showHelp {
			if key.Matches(msg, m.keys.Help) || key.Matches(msg, m.keys.Back) || key.Matches(msg, m.keys.Quit) {
				m.showHelp = false
				return m, nil
			}
			return m, nil
		}

		if m.currentView == TimerView {
			if m.timerView.IsShowingSelector() {
				return m.handleSelectorKeys(msg)
			}
			if m.timerView.GetTimerComponent().IsEditingDescription() {
				return m.handleDescriptionEditKeys(msg)
			}
		}

		switch {
		case key.Matches(msg, m.keys.Help):
			m.showHelp = true
			return m, nil

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.SwitchToTimer):
			m.currentView = TimerView
			m.statusBar.SetInfo("Switched to Timer view")
			return m, nil

		case key.Matches(msg, m.keys.SwitchToEntries):
			m.currentView = EntriesView
			m.statusBar.SetInfo("Switched to Entries view")
			return m, m.loadEntries()

		case key.Matches(msg, m.keys.SwitchToReports):
			m.currentView = ReportsView
			m.statusBar.SetInfo("Switched to Reports view")
			return m, m.loadReports()

		case key.Matches(msg, m.keys.Refresh):
			return m, m.refresh()

		case key.Matches(msg, m.keys.Left):
			if m.currentView == ReportsView {
				m.reportsView.PrevDate()
				return m, m.loadReports()
			}

		case key.Matches(msg, m.keys.Right):
			if m.currentView == ReportsView {
				m.reportsView.NextDate()
				return m, m.loadReports()
			}

		case key.Matches(msg, m.keys.Up):
			if m.currentView == EntriesView {
				m.entriesView.MoveUp()
				return m, nil
			}

		case key.Matches(msg, m.keys.Down):
			if m.currentView == EntriesView {
				m.entriesView.MoveDown()
				return m, nil
			}

		case key.Matches(msg, m.keys.ToggleView):
			switch m.currentView {
			case EntriesView:
				m.entriesView.ToggleViewMode()
				return m, m.loadEntries()
			case ReportsView:
				m.reportsView.ToggleReportType()
				return m, m.loadReports()
			}

		case key.Matches(msg, m.keys.StartTimer):
			if m.currentView == TimerView && !m.timerService.GetState().IsRunning {
				m.timerView.ShowProjectSelector()
				return m, nil
			}
			if m.currentView == EntriesView {
				return m.startTimerFromSelectedEntry()
			}

		case key.Matches(msg, m.keys.StopTimer):
			if m.currentView == TimerView && m.timerService.GetState().IsRunning {
				return m, m.stopTimer
			}

		case key.Matches(msg, m.keys.SelectProject):
			if m.currentView == TimerView {
				m.timerView.ShowProjectSelector()
				return m, nil
			}

		case key.Matches(msg, m.keys.EditDescription):
			if m.currentView == TimerView {
				if m.timerService.GetState().IsRunning {
					m.timerView.GetTimerComponent().StartEditingDescription()
					return m, nil
				} else {
					m.statusBar.SetInfo("No timer running to edit")
					return m, nil
				}
			}
		}

	case TickMsg:
		return m, tickCmd()

	case TimerStartedMsg:
		m.timerService.GetState().Start(msg.Entry)
		m.statusBar.SetSuccess("Timer started")
		return m, nil

	case TimerStoppedMsg:
		m.timerService.GetState().Stop()
		m.timerView.GetTimerComponent().ClearEditState()
		m.statusBar.SetSuccess("Timer stopped")
		return m, nil

	case TimerAlreadyStoppedMsg:
		m.timerService.GetState().Stop()
		m.timerView.GetTimerComponent().ClearEditState()
		m.statusBar.SetError(fmt.Errorf("timer was already stopped by other instance"))
		return m, nil

	case TimerDescriptionUpdatedMsg:
		m.timerService.GetState().Description = msg.Entry.Description
		m.timerService.GetState().TagIDs = msg.Entry.TagIDs
		m.statusBar.SetSuccess("Description and tags updated")
		return m, nil

	case ProjectsLoadedMsg:
		m.projects = msg.Projects
		m.timerView.SetProjects(msg.Projects)
		projectMap := make(map[string]string)
		for _, p := range msg.Projects {
			projectMap[p.ID] = p.Name
		}
		m.projectsMap = projectMap
		m.timerView.SetProjectMap(projectMap)
		m.entriesView.SetProjects(projectMap)
		return m, nil

	case TasksLoadedMsg:
		for _, task := range msg.Tasks {
			m.tasksMap[task.ID] = task.Name
		}
		m.timerView.GetProjectSelector().SetTasks(msg.Tasks)
		m.entriesView.SetTasks(m.tasksMap)
		return m, nil

	case TagsLoadedMsg:
		m.tags = msg.Tags
		tagMap := make(map[string]string)
		for _, tag := range msg.Tags {
			tagMap[tag.ID] = tag.Name
		}
		m.tagsMap = tagMap
		m.timerView.GetProjectSelector().SetTags(msg.Tags)
		m.timerView.SetTagMap(tagMap)
		m.entriesView.SetTags(tagMap)
		m.reportsView.SetTags(tagMap)
		return m, nil

	case TimeEntriesLoadedMsg:
		m.entries = msg.Entries
		m.entriesView.SetEntries(msg.Entries)
		return m, nil

	case DailyReportLoadedMsg:
		if report, ok := msg.Report.(*domain.DailySummary); ok {
			m.reportsView.SetDailyReport(report)
		}
		return m, nil

	case WeeklyReportLoadedMsg:
		if report, ok := msg.Report.(*domain.WeeklySummary); ok {
			m.reportsView.SetWeeklyReport(report)
		}
		return m, nil

	case ErrorMsg:
		m.statusBar.SetError(msg.Err)
		return m, nil
	}

	return m, nil
}

func (m App) handleDescriptionEditKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	timerComp := m.timerView.GetTimerComponent()

	switch msg.Type {
	case tea.KeyEnter:
		newDescription := timerComp.GetEditedDescription()
		timerComp.CancelEditingDescription()

		currentTagIDs := []string{}
		if m.timerService.GetState().CurrentEntry != nil {
			currentTagIDs = m.timerService.GetState().CurrentEntry.TagIDs
		}

		m.timerView.ShowTagSelectorForEditing(newDescription, currentTagIDs, m.tags)
		return m, nil

	case tea.KeyBackspace:
		timerComp.DeleteCharFromEdit()
		return m, nil

	case tea.KeyEsc:
		timerComp.CancelEditingDescription()
		return m, nil

	case tea.KeySpace:
		timerComp.AddCharToEdit(' ')
		return m, nil

	case tea.KeyRunes:
		for _, r := range msg.Runes {
			timerComp.AddCharToEdit(r)
		}
		return m, nil
	}

	return m, nil
}

func (m App) handleSelectorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	selector := m.timerView.GetProjectSelector()

	if selector.GetMode() == components.SelectingTags {
		switch {
		case key.Matches(msg, m.keys.Up):
			selector.MoveUp()
			return m, nil

		case key.Matches(msg, m.keys.Down):
			selector.MoveDown()
			return m, nil

		case key.Matches(msg, m.keys.Space):
			selector.ToggleCurrentTag()
			return m, nil

		case key.Matches(msg, m.keys.Enter):
			if m.timerView.IsEditingMode() {
				newDescription := m.timerView.GetEditedDescription()
				newTagIDs := selector.GetSelectedTagIDs()
				selector.Reset()
				m.timerView.HideProjectSelector()
				return m, m.updateTimerDescriptionAndTags(newDescription, newTagIDs)
			} else {
				projectID, taskID, description, tagIDs := selector.ConfirmTags()
				if projectID != nil {
					selector.Reset()
					m.timerView.HideProjectSelector()
					return m, m.startTimerWithTags(projectID, taskID, *description, tagIDs)
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Back):
			if m.timerView.IsEditingMode() {
				m.timerView.HideProjectSelector()
			} else {
				if selector.Back() {
					m.timerView.HideProjectSelector()
				}
			}
			return m, nil
		}
		return m, nil
	}

	if selector.GetMode() == components.EnteringDescription {
		switch msg.Type {
		case tea.KeyEnter:
			selector.TransitionToTagSelection()
			return m, nil

		case tea.KeyBackspace:
			selector.DeleteChar()
			return m, nil

		case tea.KeyEsc:
			if selector.Back() {
				m.timerView.HideProjectSelector()
			}
			return m, nil

		case tea.KeySpace:
			selector.AddChar(' ')
			return m, nil

		case tea.KeyRunes:
			for _, r := range msg.Runes {
				selector.AddChar(r)
			}
			return m, nil
		}
		return m, nil
	}

	switch {
	case key.Matches(msg, m.keys.Up):
		selector.MoveUp()
		return m, nil

	case key.Matches(msg, m.keys.Down):
		selector.MoveDown()
		return m, nil

	case key.Matches(msg, m.keys.Enter):
		projectID, _, needsTasks := selector.GetSelection()
		if needsTasks {
			return m, m.loadTasksForProject(*projectID)
		}
		return m, nil

	case key.Matches(msg, m.keys.Back):
		if selector.Back() {
			m.timerView.HideProjectSelector()
		}
		return m, nil
	}

	return m, nil
}

func (m App) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.showHelp {
		return m.renderHelp()
	}

	var content string

	tabs := m.renderTabs()
	content += tabs + "\n\n"

	switch m.currentView {
	case TimerView:
		content += m.renderTimerView()
	case EntriesView:
		content += m.renderEntriesView()
	case ReportsView:
		content += m.renderReportsView()
	}

	availableHeight := m.height - lipgloss.Height(tabs) - 3
	contentHeight := availableHeight - 2

	styledContent := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(contentHeight).
		Render(content)

	statusBar := m.statusBar.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		styledContent,
		statusBar,
	)
}

func (m App) renderTabs() string {
	tabs := []string{}

	timerTab := "Timer"
	entriesTab := "Entries"
	reportsTab := "Reports"

	switch m.currentView {
	case TimerView:
		tabs = append(tabs, ActiveTabStyle.Render(timerTab))
		tabs = append(tabs, InactiveTabStyle.Render(entriesTab))
		tabs = append(tabs, InactiveTabStyle.Render(reportsTab))
	case EntriesView:
		tabs = append(tabs, InactiveTabStyle.Render(timerTab))
		tabs = append(tabs, ActiveTabStyle.Render(entriesTab))
		tabs = append(tabs, InactiveTabStyle.Render(reportsTab))
	default:
		tabs = append(tabs, InactiveTabStyle.Render(timerTab))
		tabs = append(tabs, InactiveTabStyle.Render(entriesTab))
		tabs = append(tabs, ActiveTabStyle.Render(reportsTab))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m App) renderTimerView() string {
	return m.timerView.View()
}

func (m App) renderEntriesView() string {
	return m.entriesView.View()
}

func (m App) renderReportsView() string {
	return m.reportsView.View()
}

func (m App) renderHelp() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.PrimaryColor).
		MarginBottom(1).
		Align(lipgloss.Center)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.BlueColor).
		MarginTop(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(theme.GreenColor).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(theme.TextColor)

	helpContent := titleStyle.Render("⌨  Keyboard Shortcuts") + "\n\n"

	helpContent += sectionStyle.Render("Global") + "\n"
	helpContent += "  " + keyStyle.Render("1") + " " + descStyle.Render("Switch to Timer view") + "\n"
	helpContent += "  " + keyStyle.Render("2") + " " + descStyle.Render("Switch to Time Entries view") + "\n"
	helpContent += "  " + keyStyle.Render("3") + " " + descStyle.Render("Switch to Reports view") + "\n"
	helpContent += "  " + keyStyle.Render("r") + " " + descStyle.Render("Refresh current view") + "\n"
	helpContent += "  " + keyStyle.Render("?") + " " + descStyle.Render("Show this help screen") + "\n"
	helpContent += "  " + keyStyle.Render("q / Ctrl+C") + " " + descStyle.Render("Quit application") + "\n"

	helpContent += sectionStyle.Render("Timer View") + "\n"
	helpContent += "  " + keyStyle.Render("s") + " " + descStyle.Render("Start timer (opens project selector)") + "\n"
	helpContent += "  " + keyStyle.Render("x") + " " + descStyle.Render("Stop running timer") + "\n"
	helpContent += "  " + keyStyle.Render("p") + " " + descStyle.Render("Select project/task") + "\n"
	helpContent += "  " + keyStyle.Render("d") + " " + descStyle.Render("Edit description & tags of running timer") + "\n"

	helpContent += sectionStyle.Render("Time Entries View") + "\n"
	helpContent += "  " + keyStyle.Render("↑/↓ or k/j") + " " + descStyle.Render("Navigate entries") + "\n"
	helpContent += "  " + keyStyle.Render("t") + " " + descStyle.Render("Toggle between Today/This Week") + "\n"

	helpContent += sectionStyle.Render("Reports View") + "\n"
	helpContent += "  " + keyStyle.Render("←/→ or h/l") + " " + descStyle.Render("Navigate dates (prev/next day or week)") + "\n"
	helpContent += "  " + keyStyle.Render("t") + " " + descStyle.Render("Toggle between Daily/Weekly report") + "\n"

	helpContent += sectionStyle.Render("Project/Task/Tag Selector") + "\n"
	helpContent += "  " + keyStyle.Render("↑/↓ or k/j") + " " + descStyle.Render("Navigate list") + "\n"
	helpContent += "  " + keyStyle.Render("Space") + " " + descStyle.Render("Toggle tag selection (when selecting tags)") + "\n"
	helpContent += "  " + keyStyle.Render("Enter") + " " + descStyle.Render("Confirm selection") + "\n"
	helpContent += "  " + keyStyle.Render("Esc") + " " + descStyle.Render("Go back or cancel") + "\n"

	helpContent += "\n\n" + lipgloss.NewStyle().
		Foreground(theme.MutedColor).
		Italic(true).
		Render("Press ? or Esc to close this help screen")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.PrimaryColor).
		Padding(2, 4).
		Width(m.width - 4).
		Height(m.height - 2)

	return boxStyle.Render(helpContent)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *App) loadCurrentTimer() tea.Msg {
	entry, err := m.timerService.GetCurrentTimer()
	if err != nil {
		return ErrorMsg{Err: err}
	}

	if entry != nil {
		return TimerStartedMsg{Entry: entry}
	}

	return nil
}

func (m *App) loadProjects() tea.Msg {
	projects, err := m.projectService.GetAllProjects()
	if err != nil {
		return ErrorMsg{Err: err}
	}

	for _, project := range projects {
		tasks, err := m.projectService.GetTasksForProject(project.ID)
		if err == nil {
			for _, task := range tasks {
				m.tasksMap[task.ID] = task.Name
			}
		}
	}

	return ProjectsLoadedMsg{Projects: projects}
}

func (m *App) loadTasksForProject(projectID string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.projectService.GetTasksForProject(projectID)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return TasksLoadedMsg{
			ProjectID: projectID,
			Tasks:     tasks,
		}
	}
}

func (m *App) startTimerWithTags(projectID, taskID *string, description string, tagIDs []string) tea.Cmd {
	return func() tea.Msg {
		entry, err := m.timerService.StartTimer(description, projectID, taskID, tagIDs)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return TimerStartedMsg{Entry: entry}
	}
}

func (m *App) startTimerFromSelectedEntry() (*App, tea.Cmd) {
	selectedEntry := m.entriesView.GetSelectedEntry()
	if selectedEntry == nil {
		m.statusBar.SetError(fmt.Errorf("no entry selected"))
		return m, nil
	}

	// Extract all properties from the selected entry
	description := selectedEntry.Description
	projectID := selectedEntry.ProjectID
	taskID := selectedEntry.TaskID
	tagIDs := selectedEntry.TagIDs
	if tagIDs == nil {
		tagIDs = []string{}
	}

	m.statusBar.SetInfo("Starting timer from entry...")
	return m, m.startTimerWithTags(projectID, taskID, description, tagIDs)
}

func (m *App) stopTimer() tea.Msg {
	entry, alreadyStopped, err := m.timerService.StopTimer()
	if err != nil {
		return ErrorMsg{Err: err}
	}
	if alreadyStopped {
		return TimerAlreadyStoppedMsg{}
	}

	return TimerStoppedMsg{Entry: entry}
}

func (m *App) updateTimerDescriptionAndTags(description string, tagIDs []string) tea.Cmd {
	return func() tea.Msg {
		if m.timerService.GetState().CurrentEntry == nil {
			return ErrorMsg{Err: fmt.Errorf("no timer entry to update")}
		}

		entryID := m.timerService.GetState().CurrentEntry.ID
		currentEntry := m.timerService.GetState().CurrentEntry

		req := api.TimeEntryRequest{
			Start:       currentEntry.TimeInterval.Start,
			End:         currentEntry.TimeInterval.End,
			Description: description,
			ProjectID:   currentEntry.ProjectID,
			TaskID:      currentEntry.TaskID,
			TagIDs:      tagIDs,
		}

		entry, err := m.timerService.UpdateTimeEntry(entryID, req)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return TimerDescriptionUpdatedMsg{Entry: entry}
	}
}

func (m *App) loadEntries() tea.Cmd {
	return func() tea.Msg {
		var entries []api.TimeEntry
		var err error

		if m.entriesView.GetViewMode() == components.ViewToday {
			entries, err = m.entryService.GetEntriesForToday()
		} else {
			entries, err = m.entryService.GetEntriesForWeek()
		}

		if err != nil {
			return ErrorMsg{Err: err}
		}

		return TimeEntriesLoadedMsg{Entries: entries}
	}
}

func (m *App) loadReports() tea.Cmd {
	return func() tea.Msg {
		selectedDate := m.reportsView.GetSelectedDate()

		if m.reportsView.GetReportType() == components.DailyReport {
			report, err := m.reportService.GetDailySummary(selectedDate, m.projectsMap, m.tasksMap)
			if err != nil {
				return ErrorMsg{Err: err}
			}
			return DailyReportLoadedMsg{
				Date:   selectedDate,
				Report: report,
			}
		} else {
			weekday := int(selectedDate.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			weekStart := selectedDate.AddDate(0, 0, -(weekday - 1))

			report, err := m.reportService.GetWeeklySummary(weekStart, m.projectsMap, m.tasksMap)
			if err != nil {
				return ErrorMsg{Err: err}
			}
			return WeeklyReportLoadedMsg{
				StartDate: weekStart,
				Report:    report,
			}
		}
	}
}

func (m *App) loadTags() tea.Msg {
	tags, err := m.tagService.GetAllTags()
	if err != nil {
		return ErrorMsg{Err: err}
	}

	return TagsLoadedMsg{Tags: tags}
}

func (m *App) refresh() tea.Cmd {
	switch m.currentView {
	case EntriesView:
		return m.loadEntries()
	case ReportsView:
		return m.loadReports()
	default:
		return m.loadCurrentTimer
	}
}
