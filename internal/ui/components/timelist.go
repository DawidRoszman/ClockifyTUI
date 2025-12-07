package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"main/internal/api"
	"main/internal/domain"
	"main/internal/ui/theme"
)

type EntriesViewMode int

const (
	ViewToday EntriesViewMode = iota
	ViewThisWeek
)

type EntriesComponent struct {
	entries       []api.TimeEntry
	projects      map[string]string
	tasks         map[string]string
	tags          map[string]string
	selectedIndex int
	viewMode      EntriesViewMode
	width         int
	height        int
}

var (
	entryItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			MarginBottom(1)

	entrySelectedStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(theme.MauveColor).
				Bold(true).
				MarginBottom(1)

	entryTimeStyle = lipgloss.NewStyle().
			Foreground(theme.Subtext0Color)

	entryDurationStyle = lipgloss.NewStyle().
				Foreground(theme.GreenColor).
				Bold(true)
)

func NewEntriesComponent() *EntriesComponent {
	return &EntriesComponent{
		viewMode:      ViewToday,
		selectedIndex: 0,
		projects:      make(map[string]string),
		tasks:         make(map[string]string),
	}
}

func (c *EntriesComponent) SetEntries(entries []api.TimeEntry) {
	c.entries = entries
	if c.selectedIndex >= len(entries) {
		c.selectedIndex = 0
	}
	if len(entries) > 0 && c.selectedIndex < 0 {
		c.selectedIndex = 0
	}
}

func (c *EntriesComponent) SetProjects(projects map[string]string) {
	c.projects = projects
}

func (c *EntriesComponent) SetTasks(tasks map[string]string) {
	c.tasks = tasks
}

func (c *EntriesComponent) SetTags(tags map[string]string) {
	c.tags = tags
}

func (c *EntriesComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

func (c *EntriesComponent) SetViewMode(mode EntriesViewMode) {
	c.viewMode = mode
	c.selectedIndex = 0
}

func (c *EntriesComponent) GetViewMode() EntriesViewMode {
	return c.viewMode
}

func (c *EntriesComponent) ToggleViewMode() {
	if c.viewMode == ViewToday {
		c.viewMode = ViewThisWeek
	} else {
		c.viewMode = ViewToday
	}
	c.selectedIndex = 0
}

func (c *EntriesComponent) NextItem() {
	if c.selectedIndex < len(c.entries)-1 {
		c.selectedIndex++
	}
}

func (c *EntriesComponent) PrevItem() {
	if c.selectedIndex > 0 {
		c.selectedIndex--
	}
}

func (c *EntriesComponent) View() string {
	if len(c.entries) == 0 {
		return c.renderEmpty()
	}

	content := ""
	maxVisible := 8

	visibleStart := 0
	visibleEnd := len(c.entries)

	if len(c.entries) > maxVisible {
		if c.selectedIndex > maxVisible/2 {
			visibleStart = c.selectedIndex - maxVisible/2
		}
		visibleEnd = visibleStart + maxVisible
		if visibleEnd > len(c.entries) {
			visibleEnd = len(c.entries)
			visibleStart = max(visibleEnd - maxVisible, 0)
		}
	}

	for i := visibleStart; i < visibleEnd; i++ {
		entry := c.entries[i]
		entryView := c.formatEntry(&entry, i == c.selectedIndex)
		content += entryView + "\n"
	}

	helpText := lipgloss.NewStyle().Foreground(theme.Subtext0Color).
		Render(fmt.Sprintf("↑/↓: navigate | t: toggle view (%d entries)", len(c.entries)))

	return content + "\n" + helpText
}

func (c *EntriesComponent) renderEmpty() string {
	emptyStyle := lipgloss.NewStyle().
		Foreground(theme.Subtext0Color).
		Italic(true)

	modeStr := "today"
	if c.viewMode == ViewThisWeek {
		modeStr = "this week"
	}

	return emptyStyle.Render(fmt.Sprintf("No time entries for %s", modeStr))
}

func (c *EntriesComponent) formatEntry(entry *api.TimeEntry, selected bool) string {
	start := entry.TimeInterval.Start.Local()
	startTime := start.Format("15:04")

	var duration string
	if entry.TimeInterval.End != nil {
		end := entry.TimeInterval.End.Local()
		endTime := end.Format("15:04")
		dur := end.Sub(start)
		duration = domain.FormatDuration(dur)
		startTime = fmt.Sprintf("%s - %s", startTime, endTime)
	} else {
		duration = "Running"
		startTime = fmt.Sprintf("%s - now", startTime)
	}

	description := entry.Description
	if description == "" {
		description = "(no description)"
	}

	projectName := ""
	if entry.ProjectID != nil {
		if name, ok := c.projects[*entry.ProjectID]; ok {
			projectName = name
		} else {
			projectName = *entry.ProjectID
		}
	}

	taskName := ""
	if entry.TaskID != nil {
		if name, ok := c.tasks[*entry.TaskID]; ok {
			taskName = " • " + name
		}
	}

	tagsStr := ""
	if len(entry.TagIDs) > 0 {
		tagNames := []string{}
		for _, tagID := range entry.TagIDs {
			if name, ok := c.tags[tagID]; ok {
				tagNames = append(tagNames, name)
			}
		}
		if len(tagNames) > 0 {
			tagStyle := lipgloss.NewStyle().
				Foreground(theme.MauveColor).
				Italic(true)

			joinedTags := ""
			for i, name := range tagNames {
				if i > 0 {
					joinedTags += ", "
				}
				joinedTags += name
			}
			tagsStr = " • " + tagStyle.Render(joinedTags)
		}
	}

	line1 := description
	line2 := entryTimeStyle.Render(startTime) + " " + entryDurationStyle.Render(duration)
	if projectName != "" {
		line2 += " • " + projectName + taskName
	}
	line2 += tagsStr

	content := line1 + "\n" + line2

	if selected {
		return entrySelectedStyle.Render("▶ " + content)
	}
	return entryItemStyle.Render(content)
}
