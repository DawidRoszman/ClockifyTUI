package domain

import (
	"fmt"
	"time"

	"main/internal/api"
)

type TimerState struct {
	IsRunning    bool
	CurrentEntry *api.TimeEntry
	StartTime    time.Time
	Description  string
	ProjectID    *string
	TaskID       *string
	TagIDs       []string
}

func NewTimerState() *TimerState {
	return &TimerState{
		IsRunning: false,
	}
}

func (t *TimerState) Start(entry *api.TimeEntry) {
	t.IsRunning = true
	t.CurrentEntry = entry
	t.StartTime = entry.TimeInterval.Start
	t.Description = entry.Description
	t.ProjectID = entry.ProjectID
	t.TaskID = entry.TaskID
	t.TagIDs = entry.TagIDs
}

func (t *TimerState) Stop() {
	t.IsRunning = false
	t.CurrentEntry = nil
	t.StartTime = time.Time{}
	t.Description = ""
	t.ProjectID = nil
	t.TaskID = nil
	t.TagIDs = nil
}

func (t *TimerState) GetElapsedDuration() time.Duration {
	if !t.IsRunning {
		return 0
	}
	return time.Since(t.StartTime)
}

func (t *TimerState) UpdateFromEntry(entry *api.TimeEntry) {
	if entry == nil {
		if t.IsRunning {
			t.Stop()
		}
		return
	}

	if entry.TimeInterval.End != nil {
		t.Stop()
	} else {
		t.Start(entry)
	}
}

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
