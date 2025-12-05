package api

import (
	"fmt"
	"time"
)

func (c *Client) GetTimeEntries(start, end time.Time) ([]TimeEntry, error) {
	path := fmt.Sprintf("/workspaces/%s/user/%s/time-entries?start=%s&end=%s",
		c.workspaceID,
		c.userID,
		start.UTC().Format(time.RFC3339),
		end.UTC().Format(time.RFC3339))

	var entries []TimeEntry
	if err := c.get(path, &entries); err != nil {
		return nil, fmt.Errorf("failed to get time entries: %w", err)
	}
	return entries, nil
}

func (c *Client) GetCurrentTimer() (*TimeEntry, error) {
	path := fmt.Sprintf("/workspaces/%s/user/%s/time-entries?in-progress=true",
		c.workspaceID,
		c.userID)

	var entries []TimeEntry
	if err := c.get(path, &entries); err != nil {
		return nil, fmt.Errorf("failed to get current timer: %w", err)
	}

	if len(entries) == 0 {
		return nil, nil
	}

	return &entries[0], nil
}

func (c *Client) StartTimer(description string, projectID, taskID *string) (*TimeEntry, error) {
	req := TimeEntryRequest{
		Start:       time.Now().UTC(),
		Description: description,
		ProjectID:   projectID,
		TaskID:      taskID,
	}

	path := fmt.Sprintf("/workspaces/%s/time-entries", c.workspaceID)

	var entry TimeEntry
	if err := c.post(path, req, &entry); err != nil {
		return nil, fmt.Errorf("failed to start timer: %w", err)
	}

	return &entry, nil
}

func (c *Client) StopTimer() (*TimeEntry, error) {
	now := time.Now().UTC()
	req := TimeEntryRequest{
		End: &now,
	}

	path := fmt.Sprintf("/workspaces/%s/user/%s/time-entries", c.workspaceID, c.userID)

	var entry TimeEntry
	if err := c.patch(path, req, &entry); err != nil {
		return nil, fmt.Errorf("failed to stop timer: %w", err)
	}

	return &entry, nil
}

func (c *Client) CreateTimeEntry(entry TimeEntryRequest) (*TimeEntry, error) {
	path := fmt.Sprintf("/workspaces/%s/time-entries", c.workspaceID)

	var result TimeEntry
	if err := c.post(path, entry, &result); err != nil {
		return nil, fmt.Errorf("failed to create time entry: %w", err)
	}

	return &result, nil
}

func (c *Client) UpdateTimeEntry(id string, entry TimeEntryRequest) (*TimeEntry, error) {
	path := fmt.Sprintf("/workspaces/%s/time-entries/%s", c.workspaceID, id)

	var result TimeEntry
	if err := c.patch(path, entry, &result); err != nil {
		return nil, fmt.Errorf("failed to update time entry: %w", err)
	}

	return &result, nil
}

func (c *Client) DeleteTimeEntry(id string) error {
	path := fmt.Sprintf("/workspaces/%s/time-entries/%s", c.workspaceID, id)

	if err := c.delete(path); err != nil {
		return fmt.Errorf("failed to delete time entry: %w", err)
	}

	return nil
}
