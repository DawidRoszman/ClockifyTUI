package api

import "time"

type User struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	ActiveWorkspace string `json:"activeWorkspace"`
}

type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TimeInterval struct {
	Start    time.Time  `json:"start"`
	End      *time.Time `json:"end,omitempty"`
	Duration *string    `json:"duration,omitempty"`
}

type TimeEntry struct {
	ID           string       `json:"id"`
	Description  string       `json:"description"`
	ProjectID    *string      `json:"projectId,omitempty"`
	TaskID       *string      `json:"taskId,omitempty"`
	TagIDs       []string     `json:"tagIds,omitempty"`
	WorkspaceID  string       `json:"workspaceId"`
	UserID       string       `json:"userId"`
	TimeInterval TimeInterval `json:"timeInterval"`
}

type TimeEntryRequest struct {
	Start       time.Time `json:"start"`
	End         *time.Time `json:"end,omitempty"`
	Description string    `json:"description"`
	ProjectID   *string   `json:"projectId,omitempty"`
	TaskID      *string   `json:"taskId,omitempty"`
	TagIDs      []string  `json:"tagIds,omitempty"`
}

type Project struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ClientID *string `json:"clientId,omitempty"`
	Color    string  `json:"color"`
	Archived bool    `json:"archived"`
}

type Task struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
	Status    string `json:"status"`
}

type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
	Archived    bool   `json:"archived"`
}

type DetailedReportRequest struct {
	DateRangeStart time.Time      `json:"dateRangeStart"`
	DateRangeEnd   time.Time      `json:"dateRangeEnd"`
	DetailedFilter DetailedFilter `json:"detailedFilter"`
}

type DetailedFilter struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type SummaryReportRequest struct {
	DateRangeStart time.Time     `json:"dateRangeStart"`
	DateRangeEnd   time.Time     `json:"dateRangeEnd"`
	SummaryFilter  SummaryFilter `json:"summaryFilter"`
}

type SummaryFilter struct {
	Groups []string `json:"groups"`
}

type DetailedReport struct {
	TimeEntries []ReportTimeEntry `json:"timeentries"`
	TotalsMap   []TotalMap        `json:"totals"`
}

type ReportTimeEntry struct {
	ID          string       `json:"_id"`
	Description string       `json:"description"`
	ProjectID   string       `json:"projectId"`
	ProjectName string       `json:"projectName"`
	TaskID      string       `json:"taskId"`
	TaskName    string       `json:"taskName"`
	TimeInterval TimeInterval `json:"timeInterval"`
}

type TotalMap struct {
	TotalTime int64 `json:"totalTime"`
}

type SummaryReport struct {
	GroupOne []SummaryGroup `json:"groupOne"`
}

type SummaryGroup struct {
	Duration int64           `json:"duration"`
	Name     string          `json:"name"`
	Children []SummaryGroup  `json:"children,omitempty"`
}
