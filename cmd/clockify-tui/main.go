package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"main/internal/api"
	"main/internal/config"
	"main/internal/ui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client := api.NewClient(cfg.APIKey, cfg.BaseURL)

	user, err := client.GetCurrentUser()
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	client.SetUserID(user.ID)

	workspaceID := cfg.WorkspaceID
	if workspaceID == "" {
		workspaceID = user.ActiveWorkspace
	}
	client.SetWorkspace(workspaceID)

	app := ui.NewApp(client)
	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
