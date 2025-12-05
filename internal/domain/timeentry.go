package domain

import (
	"time"

	"main/internal/api"
)

type TimeEntryService struct {
	apiClient *api.Client
}

func NewTimeEntryService(client *api.Client) *TimeEntryService {
	return &TimeEntryService{
		apiClient: client,
	}
}

func (s *TimeEntryService) GetEntriesForToday() ([]api.TimeEntry, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)
	return s.apiClient.GetTimeEntries(start, end)
}

func (s *TimeEntryService) GetEntriesForWeek() ([]api.TimeEntry, error) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
		AddDate(0, 0, -(weekday - 1))
	end := start.AddDate(0, 0, 7)
	return s.apiClient.GetTimeEntries(start, end)
}

func (s *TimeEntryService) GetEntriesForRange(start, end time.Time) ([]api.TimeEntry, error) {
	return s.apiClient.GetTimeEntries(start, end)
}

func (s *TimeEntryService) GetDurationForEntry(entry *api.TimeEntry) time.Duration {
	if entry.TimeInterval.End != nil {
		return entry.TimeInterval.End.Sub(entry.TimeInterval.Start)
	}
	return time.Since(entry.TimeInterval.Start)
}
