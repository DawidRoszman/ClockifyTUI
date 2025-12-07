package api

import "fmt"

func (c *Client) GetTasksForProject(projectID string) ([]Task, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s/tasks", c.workspaceID, projectID)

	var tasks []Task
	if err := c.get(path, &tasks); err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}
