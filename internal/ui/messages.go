package ui

import (
	"time"

	"main/internal/api"
)

type ViewType int

const (
	TimerView ViewType = iota
	EntriesView
	ReportsView
)

type TimerStartedMsg struct {
	Entry *api.TimeEntry
}

type TimerStoppedMsg struct {
	Entry *api.TimeEntry
}

type TimerAlreadyStoppedMsg struct {
}

type TimeEntriesLoadedMsg struct {
	Entries []api.TimeEntry
}

type ProjectsLoadedMsg struct {
	Projects []api.Project
}

type TasksLoadedMsg struct {
	ProjectID string
	Tasks     []api.Task
}

type TagsLoadedMsg struct {
	Tags []api.Tag
}

type DailyReportLoadedMsg struct {
	Date   time.Time
	Report any
}

type WeeklyReportLoadedMsg struct {
	StartDate time.Time
	Report    any
}

type ErrorMsg struct {
	Err error
}

type SwitchViewMsg struct {
	View ViewType
}

type TickMsg time.Time

type ProjectSelectedMsg struct {
	ProjectID *string
	TaskID    *string
}

type InitializedMsg struct {
	WorkspaceID string
	UserID      string
}

type RefreshMsg struct{}

type TimerDescriptionUpdatedMsg struct {
	Entry *api.TimeEntry
}

type DescriptionSuggestionsLoadedMsg struct {
	Suggestions []string
}
