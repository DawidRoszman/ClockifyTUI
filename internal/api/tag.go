package api

import "fmt"

func (c *Client) GetTags() ([]Tag, error) {
	path := fmt.Sprintf("/workspaces/%s/tags", c.workspaceID)

	var tags []Tag
	if err := c.get(path, &tags); err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return tags, nil
}
