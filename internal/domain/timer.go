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

type TimerService struct {
	apiClient *api.Client
	state     *TimerState
}

func NewTimerService(client *api.Client, state *TimerState) *TimerService {
	return &TimerService{
		apiClient: client,
		state:     state,
	}
}

func (s *TimerService) GetCurrentTimer() (*api.TimeEntry, error) {
	return s.apiClient.GetCurrentTimer()
}

func (s *TimerService) StartTimer(description string, projectID, taskID *string, tagIDs []string) (*api.TimeEntry, error) {
	entry, err := s.apiClient.StartTimer(description, projectID, taskID, tagIDs)
	if err != nil {
		return nil, err
	}
	s.state.Start(entry)
	return entry, nil
}

func (s *TimerService) StopTimer() (*api.TimeEntry, bool, error) {
	currentEntry, err := s.apiClient.GetCurrentTimer()
	if err != nil {
		return nil, false, err
	}
	if currentEntry == nil {
		s.state.Stop()
		return nil, true, nil
	}

	entry, err := s.apiClient.StopTimer()
	if err != nil {
		return nil, false, err
	}

	s.state.Stop()
	return entry, false, nil
}

func (s *TimerService) UpdateTimeEntry(entryID string, req api.TimeEntryRequest) (*api.TimeEntry, error) {
	entry, err := s.apiClient.UpdateTimeEntry(entryID, req)
	if err != nil {
		return nil, err
	}

	s.state.Description = entry.Description
	s.state.TagIDs = entry.TagIDs
	if s.state.CurrentEntry != nil {
		s.state.CurrentEntry.Description = entry.Description
		s.state.CurrentEntry.TagIDs = entry.TagIDs
	}

	return entry, nil
}

func (s *TimerService) GetState() *TimerState {
	return s.state
}
