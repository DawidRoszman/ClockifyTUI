 Clockify TUI Application - Implementation Plan

 Overview

 Build a terminal user interface (TUI) application for Clockify time tracking with timer controls, time entry management, and reporting capabilities.

 Technology Stack

 - Language: Go 1.25.4
 - TUI Framework: Bubbletea (The Elm Architecture)
 - Styling: Lip Gloss
 - UI Components: Bubbles
 - Authentication: Environment variable (CLOCKIFY_API_KEY)
 - API: Clockify REST API v1

 Project Structure

 clockify-tui/
 ├── cmd/clockify-tui/
 │   └── main.go                    # Application entry point
 ├── internal/
 │   ├── api/                       # Clockify API client
 │   │   ├── client.go              # HTTP client wrapper
 │   │   ├── auth.go                # Authentication
 │   │   ├── models.go              # API data structures
 │   │   ├── timeentries.go         # Time entries endpoints
 │   │   ├── projects.go            # Projects endpoints
 │   │   ├── tasks.go               # Tasks endpoints
 │   │   └── reports.go             # Reports endpoints
 │   ├── domain/                    # Business logic
 │   │   ├── timer.go               # Timer state management
 │   │   ├── timeentry.go           # Time entry operations
 │   │   ├── project.go             # Project/task management
 │   │   └── report.go              # Report calculations
 │   ├── config/
 │   │   └── config.go              # Configuration loading
 │   ├── cache/
 │   │   └── cache.go               # In-memory cache for projects/tasks
 │   └── ui/                        # Bubbletea UI layer
 │       ├── app.go                 # Main app model (Init/Update/View)
 │       ├── messages.go            # Message types
 │       ├── keys.go                # Key bindings
 │       ├── styles.go              # Lip Gloss styles
 │       ├── components/
 │       │   ├── timer.go           # Timer display component
 │       │   ├── timelist.go        # Time entries list
 │       │   ├── reports.go         # Reports display
 │       │   ├── projectselector.go # Project/task selector
 │       │   └── statusbar.go       # Status bar
 │       └── views/
 │           ├── timer.go           # Timer view
 │           ├── entries.go         # Time entries view
 │           └── reports.go         # Reports view

 Core Features (MVP)

 1. Timer Controls

 - Start/stop timer with description
 - Select project and task for time entry
 - Real-time elapsed time display
 - Current timer status indicator

 2. Time Entries View

 - Display today's time entries
 - Display current week's entries
 - Navigate entries with keyboard
 - Show project/task for each entry
 - Refresh functionality

 3. Reports View

 - Daily Summary: Hours worked today by project/task
 - Weekly Summary: Hours for current week with project breakdown
 - Task Totals: Total hours per task over selected time period
 - Date navigation (previous/next)
 - Toggle between report types

 Implementation Phases

 Phase 1: Foundation & API Client (Days 1-4)

 Goal: Set up project structure and complete API integration

 1. Create directory structure
 2. Install dependencies:
   - github.com/charmbracelet/bubbletea
   - github.com/charmbracelet/lipgloss
   - github.com/charmbracelet/bubbles
 3. Implement configuration loading from CLOCKIFY_API_KEY env var
 4. Build API client with authentication
 5. Implement all API endpoint wrappers:
   - User/workspace authentication
   - Time entries (get, start, stop, create, update, delete)
   - Projects and tasks (get, list)
   - Reports (detailed, summary)
 6. Implement in-memory cache for projects/tasks (5min TTL)
 7. Create API models matching Clockify schema

 Critical Files:
 - /cmd/clockify-tui/main.go
 - /internal/config/config.go
 - /internal/api/client.go
 - /internal/api/auth.go
 - /internal/api/models.go
 - /internal/api/timeentries.go
 - /internal/api/projects.go
 - /internal/api/tasks.go
 - /internal/api/reports.go
 - /internal/cache/cache.go

 Validation: Successfully authenticate and fetch data from Clockify API

 Phase 2: Domain Logic (Days 5-6)

 Goal: Implement business logic independent of API and UI

 1. Implement timer state management:
   - Track running/stopped state
   - Calculate elapsed duration
   - Manage current time entry
 2. Create time entry service:
   - Fetch entries for date ranges
   - Format durations
   - Filter by today/week
 3. Build project service with caching
 4. Implement report service:
   - Aggregate time entries by project/task
   - Calculate daily summaries
   - Calculate weekly summaries with trends
   - Generate task totals

 Critical Files:
 - /internal/domain/timer.go
 - /internal/domain/timeentry.go
 - /internal/domain/project.go
 - /internal/domain/report.go

 Validation: Business logic works with mocked API client

 Phase 3: Bubbletea UI Framework (Days 7-8)

 Goal: Set up TUI foundation with navigation

 1. Create main app model with Bubbletea's Init/Update/View pattern
 2. Define all message types (timer started/stopped, entries loaded, etc.)
 3. Implement key bindings:
   - 1/2/3: Switch views (Timer/Entries/Reports)
   - s: Start timer
   - x: Stop timer
   - r: Refresh
   - q: Quit
 4. Set up Lip Gloss styles (colors, borders, layout)
 5. Implement view switching with tab navigation
 6. Create status bar component

 Critical Files:
 - /internal/ui/app.go
 - /internal/ui/messages.go
 - /internal/ui/keys.go
 - /internal/ui/styles.go
 - /internal/ui/components/statusbar.go

 Validation: Can launch TUI, switch between empty views, quit properly

 Phase 4: Timer View (Days 9-10)

 Goal: Complete timer functionality

 1. Build timer component:
   - Display elapsed time (updates every second)
   - Show timer status (running/stopped)
   - Display current description, project, task
 2. Create project/task selector component:
   - List all projects
   - Navigate with arrow keys
   - Select project to show tasks
   - Confirm selection
 3. Wire up timer commands:
   - Start timer (opens project selector if needed)
   - Stop timer
   - Update UI in real-time
 4. Implement tick messages for live updates
 5. Add error handling and status feedback

 Critical Files:
 - /internal/ui/components/timer.go
 - /internal/ui/components/projectselector.go
 - /internal/ui/views/timer.go

 Validation: Can start/stop timer, select projects/tasks, see real-time updates

 Phase 5: Time Entries View (Days 11-12)

 Goal: View and manage time entries

 1. Build time entries list component:
   - Display entries with description, project, task, duration
   - Format timestamps (e.g., "9:00 AM - 10:30 AM")
   - Highlight selected entry
   - Navigate with arrow keys
 2. Implement view mode toggle (today/this week)
 3. Add refresh functionality
 4. Show loading state while fetching
 5. Handle empty states (no entries)

 Critical Files:
 - /internal/ui/components/timelist.go
 - /internal/ui/views/entries.go

 Validation: Can view and navigate time entries for today and current week

 Phase 6: Reports View (Days 13-14)

 Goal: Display summary reports

 1. Build reports component:
   - Render daily summary (hours by project/task)
   - Render weekly summary (daily breakdown, project totals)
   - Render task totals view
   - Use ASCII tables/bars for visual representation
 2. Implement date navigation:
   - Previous/next day (for daily report)
   - Previous/next week (for weekly report)
 3. Add report type toggle (daily/weekly/task totals)
 4. Format durations clearly (e.g., "8h 30m")
 5. Show project colors/indicators

 Critical Files:
 - /internal/ui/components/reports.go
 - /internal/ui/views/reports.go

 Validation: Can view all report types and navigate dates

 Phase 7: Polish & Testing (Days 15-16)

 Goal: Production-ready application

 1. Add loading spinners for async operations
 2. Improve error messages and recovery
 3. Add help screen (? key) showing all keyboard shortcuts
 4. Optimize API calls and caching
 5. Comprehensive testing:
   - Unit tests for domain logic
   - Integration tests for API
   - Manual TUI testing
 6. Write README with setup instructions
 7. Code cleanup and refactoring per SOLID principles
 8. Edge case handling (network errors, expired tokens, etc.)

 Validation: Application is stable, well-tested, and documented

 Architecture Principles

 Layered Architecture

 1. API Layer (internal/api/): HTTP communication only
 2. Domain Layer (internal/domain/): Business logic only
 3. UI Layer (internal/ui/): Presentation only
 4. Clear separation: Each layer has single responsibility

 Bubbletea Pattern

 - Model: App state and components
 - Update: Handle messages, return commands
 - View: Render current state as string
 - Messages: All state changes via message passing
 - Commands: Async operations (API calls) return messages

 SOLID Compliance

 - Single Responsibility: Each file/package has one purpose
 - Open/Closed: Message-based extension without modification
 - Liskov Substitution: Mock-friendly interfaces
 - Interface Segregation: Minimal, focused APIs
 - Dependency Inversion: Domain depends on abstractions

 Key Design Decisions

 1. Single Workspace: Use active workspace from user profile (simplifies MVP)
 2. In-Memory Cache: 5-minute TTL for projects/tasks only (not time entries)
 3. Reports API: Use Clockify's aggregated reports endpoints (not raw entries)
 4. Error Handling: Display in status bar, don't crash application
 5. No Persistence: API is source of truth, fetch fresh on startup

 Navigation & Controls

 Global Keys

 - 1: Timer view
 - 2: Time Entries view
 - 3: Reports view
 - r: Refresh current view
 - q / Ctrl+C: Quit
 - ?: Show help

 Timer View

 - s: Start timer (opens project selector)
 - x: Stop timer
 - p: Select project/task

 Entries View

 - ↑/↓: Navigate entries
 - t: Toggle today/week view

 Reports View

 - ←/→: Navigate dates/weeks
 - t: Toggle report type (daily/weekly/task totals)

 Clockify API Integration

 Base URL

 https://api.clockify.me/api/v1

 Authentication

 Header: X-Api-Key: {CLOCKIFY_API_KEY}

 Required Endpoints

 - GET /user - Get current user
 - GET /workspaces - Get workspaces
 - GET /workspaces/{id}/user/{userId}/time-entries - Get entries
 - POST /workspaces/{id}/time-entries - Start timer
 - PATCH /workspaces/{id}/user/{userId}/time-entries - Stop timer
 - GET /workspaces/{id}/projects - Get projects
 - GET /workspaces/{id}/projects/{id}/tasks - Get tasks
 - POST /workspaces/{id}/reports/summary - Summary reports
 - POST /workspaces/{id}/reports/detailed - Detailed reports

 Dependencies

 require (
     github.com/charmbracelet/bubbletea v0.25.0
     github.com/charmbracelet/lipgloss v0.9.1
     github.com/charmbracelet/bubbles v0.17.1
 )

 Build & Run

 # Set API key
 export CLOCKIFY_API_KEY="your-api-key"

 # Run
 go run cmd/clockify-tui/main.go

 # Build
 go build -o clockify-tui cmd/clockify-tui/main.go

 # Test
 go test ./...

 Success Criteria

 - User can start/stop timer with project/task selection
 - Real-time timer updates every second
 - View today's and week's time entries
 - View daily summary reports
 - View weekly summary reports
 - View total hours per task
 - Navigate with keyboard only
 - Error handling doesn't crash app
 - Clean, maintainable code following SOLID principles

 Future Enhancements (Post-MVP)

 - Edit time entries
 - Delete time entries
 - Create manual time entries
 - Tag support
 - Multi-workspace support
 - Export reports
 - Color themes
 - Custom key bindings
