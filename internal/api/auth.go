package api

import "fmt"

func (c *Client) GetCurrentUser() (*User, error) {
	var user User
	if err := c.get("/user", &user); err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	return &user, nil
}

func (c *Client) GetWorkspaces() ([]Workspace, error) {
	var workspaces []Workspace
	if err := c.get("/workspaces", &workspaces); err != nil {
		return nil, fmt.Errorf("failed to get workspaces: %w", err)
	}
	return workspaces, nil
}

func (c *Client) ValidateAPIKey() error {
	_, err := c.GetCurrentUser()
	return err
}
