package api

import "fmt"

func (c *Client) GetProjects() ([]Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects?archived=false", c.workspaceID)

	var projects []Project
	if err := c.get(path, &projects); err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	return projects, nil
}

func (c *Client) GetProjectByID(id string) (*Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s", c.workspaceID, id)

	var project Project
	if err := c.get(path, &project); err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	return &project, nil
}
