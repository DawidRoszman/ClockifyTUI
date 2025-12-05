package domain

import (
	"time"

	"main/internal/api"
)

type ReportService struct {
	apiClient *api.Client
}

type DailySummary struct {
	Date          time.Time
	TotalDuration time.Duration
	ByProject     map[string]*ProjectSummary
}

type ProjectSummary struct {
	ProjectID     string
	ProjectName   string
	TotalDuration time.Duration
	ByTask        map[string]*TaskSummary
}

type TaskSummary struct {
	TaskID   string
	TaskName string
	Duration time.Duration
}

type WeeklySummary struct {
	StartDate     time.Time
	EndDate       time.Time
	TotalDuration time.Duration
	ByDay         map[string]time.Duration
	ByProject     map[string]*ProjectSummary
}

func NewReportService(client *api.Client) *ReportService {
	return &ReportService{
		apiClient: client,
	}
}

func (s *ReportService) GetDailySummary(date time.Time, projectMap, taskMap map[string]string) (*DailySummary, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	entries, err := s.apiClient.GetTimeEntries(start, end)
	if err != nil {
		return nil, err
	}

	return s.aggregateDailySummary(date, entries, projectMap, taskMap), nil
}

func (s *ReportService) GetWeeklySummary(weekStart time.Time, projectMap, taskMap map[string]string) (*WeeklySummary, error) {
	start := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
	end := start.AddDate(0, 0, 7)

	entries, err := s.apiClient.GetTimeEntries(start, end)
	if err != nil {
		return nil, err
	}

	return s.aggregateWeeklySummary(start, end, entries, projectMap, taskMap), nil
}

func (s *ReportService) aggregateDailySummary(date time.Time, entries []api.TimeEntry, projectMap, taskMap map[string]string) *DailySummary {
	summary := &DailySummary{
		Date:      date,
		ByProject: make(map[string]*ProjectSummary),
	}

	for _, entry := range entries {
		duration := s.calculateDuration(&entry)
		summary.TotalDuration += duration

		projectID := "no-project"
		projectName := "No Project"
		if entry.ProjectID != nil {
			projectID = *entry.ProjectID
			if name, ok := projectMap[projectID]; ok {
				projectName = name
			} else {
				projectName = projectID
			}
		}

		if _, ok := summary.ByProject[projectID]; !ok {
			summary.ByProject[projectID] = &ProjectSummary{
				ProjectID:   projectID,
				ProjectName: projectName,
				ByTask:      make(map[string]*TaskSummary),
			}
		}

		summary.ByProject[projectID].TotalDuration += duration

		taskID := "no-task"
		taskName := "No Task"
		if entry.TaskID != nil {
			taskID = *entry.TaskID
			if name, ok := taskMap[taskID]; ok {
				taskName = name
			} else {
				taskName = taskID
			}
		}

		if _, ok := summary.ByProject[projectID].ByTask[taskID]; !ok {
			summary.ByProject[projectID].ByTask[taskID] = &TaskSummary{
				TaskID:   taskID,
				TaskName: taskName,
			}
		}

		summary.ByProject[projectID].ByTask[taskID].Duration += duration
	}

	return summary
}

func (s *ReportService) aggregateWeeklySummary(start, end time.Time, entries []api.TimeEntry, projectMap, taskMap map[string]string) *WeeklySummary {
	summary := &WeeklySummary{
		StartDate: start,
		EndDate:   end,
		ByDay:     make(map[string]time.Duration),
		ByProject: make(map[string]*ProjectSummary),
	}

	for _, entry := range entries {
		duration := s.calculateDuration(&entry)
		summary.TotalDuration += duration

		dayKey := entry.TimeInterval.Start.Format("2006-01-02")
		summary.ByDay[dayKey] += duration

		projectID := "no-project"
		projectName := "No Project"
		if entry.ProjectID != nil {
			projectID = *entry.ProjectID
			if name, ok := projectMap[projectID]; ok {
				projectName = name
			} else {
				projectName = projectID
			}
		}

		if _, ok := summary.ByProject[projectID]; !ok {
			summary.ByProject[projectID] = &ProjectSummary{
				ProjectID:   projectID,
				ProjectName: projectName,
				ByTask:      make(map[string]*TaskSummary),
			}
		}

		summary.ByProject[projectID].TotalDuration += duration

		taskID := "no-task"
		taskName := "No Task"
		if entry.TaskID != nil {
			taskID = *entry.TaskID
			if name, ok := taskMap[taskID]; ok {
				taskName = name
			} else {
				taskName = taskID
			}
		}

		if _, ok := summary.ByProject[projectID].ByTask[taskID]; !ok {
			summary.ByProject[projectID].ByTask[taskID] = &TaskSummary{
				TaskID:   taskID,
				TaskName: taskName,
			}
		}

		summary.ByProject[projectID].ByTask[taskID].Duration += duration
	}

	return summary
}

func (s *ReportService) calculateDuration(entry *api.TimeEntry) time.Duration {
	if entry.TimeInterval.End != nil {
		return entry.TimeInterval.End.Sub(entry.TimeInterval.Start)
	}
	return time.Since(entry.TimeInterval.Start)
}
