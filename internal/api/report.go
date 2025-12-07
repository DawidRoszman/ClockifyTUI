package api

import (
	"fmt"
	"time"
)

func (c *Client) GetDetailedReport(start, end time.Time) (*DetailedReport, error) {
	req := DetailedReportRequest{
		DateRangeStart: start,
		DateRangeEnd:   end,
		DetailedFilter: DetailedFilter{
			Page:     1,
			PageSize: 1000,
		},
	}

	path := fmt.Sprintf("/workspaces/%s/reports/detailed", c.workspaceID)

	var report DetailedReport
	if err := c.post(path, req, &report); err != nil {
		return nil, fmt.Errorf("failed to get detailed report: %w", err)
	}

	return &report, nil
}

func (c *Client) GetSummaryReport(start, end time.Time, groups []string) (*SummaryReport, error) {
	if groups == nil {
		groups = []string{"PROJECT", "TASK"}
	}

	req := SummaryReportRequest{
		DateRangeStart: start,
		DateRangeEnd:   end,
		SummaryFilter: SummaryFilter{
			Groups: groups,
		},
	}

	path := fmt.Sprintf("/workspaces/%s/reports/summary", c.workspaceID)

	var report SummaryReport
	if err := c.post(path, req, &report); err != nil {
		return nil, fmt.Errorf("failed to get summary report: %w", err)
	}

	return &report, nil
}
