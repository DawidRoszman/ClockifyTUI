package config

import (
	"errors"
	"os"
)

type Config struct {
	APIKey      string
	WorkspaceID string
	BaseURL     string
}

func Load() (*Config, error) {
	apiKey := os.Getenv("CLOCKIFY_API_KEY")
	if apiKey == "" {
		return nil, errors.New("CLOCKIFY_API_KEY environment variable is required")
	}

	workspaceID := os.Getenv("CLOCKIFY_WORKSPACE_ID")

	baseURL := os.Getenv("CLOCKIFY_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.clockify.me/api/v1"
	}

	cfg := &Config{
		APIKey:      apiKey,
		WorkspaceID: workspaceID,
		BaseURL:     baseURL,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("API key is required")
	}
	if c.BaseURL == "" {
		return errors.New("base URL is required")
	}
	return nil
}
