package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient  *http.Client
	baseURL     string
	apiKey      string
	workspaceID string
	userID      string
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (c *Client) SetWorkspace(workspaceID string) {
	c.workspaceID = workspaceID
}

func (c *Client) SetUserID(userID string) {
	c.userID = userID
}

func (c *Client) GetWorkspaceID() string {
	return c.workspaceID
}

func (c *Client) GetUserID() string {
	return c.userID
}

func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *Client) get(path string, result interface{}) error {
	return c.doRequest("GET", path, nil, result)
}

func (c *Client) post(path string, body interface{}, result interface{}) error {
	return c.doRequest("POST", path, body, result)
}

func (c *Client) patch(path string, body interface{}, result interface{}) error {
	return c.doRequest("PATCH", path, body, result)
}

func (c *Client) delete(path string) error {
	return c.doRequest("DELETE", path, nil, nil)
}
