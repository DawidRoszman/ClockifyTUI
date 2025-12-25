# Clockify TUI

A terminal user interface (TUI) for [Clockify](https://clockify.me) time tracking, built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea).

<img width="898" height="548" alt="image" src="https://github.com/user-attachments/assets/02cf4596-09f0-4d8e-8b23-a6c649b5066a" />

## Features

- ⏱️  **Timer Management**: Start/stop timers with project and task selection
- 📋 **Time Entries**: View today's and this week's time entries
- 📊 **Reports**: Daily and weekly summaries with project/task breakdowns
- ⌨️  **Keyboard-Driven**: Full keyboard navigation and control
- 🎨 **Beautiful UI**: Clean, colorful interface inspired by the Clockify web app
- ⚡ **Fast & Efficient**: In-memory caching for quick project/task lookups

## Screenshots

### Timer View
Start and stop timers with real-time duration tracking.

<img width="805" height="531" alt="image" src="https://github.com/user-attachments/assets/b4d7dac0-9c55-45b6-a0ad-037afd127f72" />

### Time Entries View
Browse through your time entries for today or the current week.

<img width="795" height="792" alt="image" src="https://github.com/user-attachments/assets/c73b0fc2-e425-436c-ae42-7243c9c6ca74" />

### Reports View
View daily and weekly summaries with visual breakdowns.

<img width="795" height="673" alt="image" src="https://github.com/user-attachments/assets/2dafadaf-561c-4566-b4f5-bb86af9c9e78" />

## Installation

### Prerequisites

- Go 1.25.4 or higher
- A Clockify account with API access

### Build from Source

```bash
# Clone the repository
git clone <url> clockify-tui
cd clockify-tui

# Build the application
go build -o clockify-tui cmd/clockify-tui/main.go

# Run the application
export CLOCKIFY_API_KEY="your-api-key-here"
./clockify-tui
```

## Configuration

The application requires a Clockify API key to authenticate.

### Getting Your API Key

1. Log in to [Clockify](https://clockify.me)
2. Go to Preferences → Advanced
3. Manage API keys
4. Copy the API key

### Environment Variables

- **`CLOCKIFY_API_KEY`** (required): Your Clockify API key
- **`CLOCKIFY_WORKSPACE_ID`** (optional): Specific workspace ID (defaults to active workspace)
- **`CLOCKIFY_BASE_URL`** (optional): Custom API base URL (defaults to `https://api.clockify.me/api/v1`)

### Example

```bash
export CLOCKIFY_API_KEY="your-api-key-here"
./clockify-tui
```

## Usage

### Keyboard Shortcuts

#### Global
- `1` - Switch to Timer view
- `2` - Switch to Time Entries view
- `3` - Switch to Reports view
- `r` - Refresh current view
- `?` - Show help screen
- `q` or `Ctrl+C` - Quit application

#### Timer View
- `s` - Start timer (opens project selector)
- `x` - Stop running timer
- `p` - Select project/task for timer

#### Time Entries View
- `↑/↓` or `k/j` - Navigate entries
- `t` - Toggle between Today/This Week

#### Reports View
- `←/→` or `h/l` - Navigate dates (previous/next day or week)
- `t` - Toggle between Daily/Weekly report

#### Project/Task Selector
- `↑/↓` or `k/j` - Navigate list
- `Enter` - Select item
- `Esc` - Go back or cancel

## Architecture

The application follows clean architecture principles with clear separation of concerns:

```
clockify-tui/
├── cmd/clockify-tui/     # Application entry point
├── internal/
│   ├── api/              # Clockify API client
│   ├── domain/           # Business logic
│   ├── config/           # Configuration management
│   ├── cache/            # In-memory caching
│   └── ui/               # Bubbletea UI components
│       ├── components/   # Reusable UI components
│       └── views/        # View implementations
```

### Design Principles

- **SOLID Principles**: Each component has a single responsibility
- **The Elm Architecture**: Bubbletea's Model-Update-View pattern
- **Clean Architecture**: Clear separation between API, domain, and UI layers
- **Dependency Inversion**: Domain logic independent of external services

## Development

### Project Structure

- **API Layer** (`internal/api/`): HTTP communication with Clockify API
- **Domain Layer** (`internal/domain/`): Business logic and data transformations
- **UI Layer** (`internal/ui/`): Bubbletea components and views
- **Configuration** (`internal/config/`): Environment variable loading
- **Cache** (`internal/cache/`): In-memory caching with TTL

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Development build
go build -o clockify-tui cmd/clockify-tui/main.go

# Production build with optimizations
go build -ldflags="-s -w" -o clockify-tui cmd/clockify-tui/main.go
```

## Features in Detail

### Timer Management
- Real-time elapsed time display (updates every second)
- Project and task selection via keyboard-navigable selector
- Visual indication of running/stopped state
- Seamless start/stop operations

### Time Entries
- View mode toggle: Today or This Week
- Displays entry description, time range, duration
- Shows associated project and task names
- Keyboard navigation through entries
- Handles empty states gracefully
- Start new entry using currently focused entry using 's'

### Reports
- **Daily Reports**: Hours by project and task for a specific day
- **Weekly Reports**: Daily breakdown with visual bars, total hours by project
- Date navigation to view historical data
- Visual bars showing relative time distribution
- Sorted by duration (most time first)

### Caching
- Projects and tasks are cached for 5 minutes
- Reduces API calls and improves responsiveness
- Automatic refresh on cache expiration
- Manual refresh available via `r` key

## Troubleshooting

### Authentication Errors

If you see authentication errors:
1. Verify your API key is correct
2. Ensure the `CLOCKIFY_API_KEY` environment variable is set
3. Check that your Clockify account is active

### No Projects Showing

If projects aren't appearing:
1. Verify you have projects in your Clockify workspace
2. Try pressing `r` to refresh
3. Check that projects aren't archived in Clockify

### Performance Issues

If the app feels slow:
1. Check your internet connection
2. The cache may have expired - it will refresh automatically
3. Try refreshing manually with `r`

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is open source and available under the MIT License.

## Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea)
- Styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- UI components from [Bubbles](https://github.com/charmbracelet/bubbles)
- Powered by the [Clockify API](https://clockify.me/developers-api)

## Support

For issues, questions, or suggestions:
- Open an issue on GitHub
- Check the built-in help screen (press `?`)

---

*This application provides a text-based interface for interacting with an external service through its API. By using this tool, you acknowledge that any actions you perform may affect the behavior, configuration, or data of that service. Use this application at your own risk. The developer is not responsible for any unintended consequences, errors, data loss, or service disruptions resulting from its use.*
