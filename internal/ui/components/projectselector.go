package components

import (
	"github.com/charmbracelet/lipgloss"
	"main/internal/api"
	"main/internal/ui/theme"
)

type SelectorMode int

const (
	SelectingProject SelectorMode = iota
	SelectingTask
	EnteringDescription
	SelectingTags
)

type ProjectSelectorComponent struct {
	projects        []api.Project
	tasks           []api.Task
	tags            []api.Tag
	selectedProject int
	selectedTask    int
	selectedTags    map[int]bool
	currentTagCursor int
	mode            SelectorMode
	filterInput     string
	description     string
	width           int
	height          int
}

var (
	selectorTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(theme.LavenderColor).
				Padding(0, 1)

	selectorItemStyle = lipgloss.NewStyle().
				PaddingLeft(2)

	selectorSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(theme.MauveColor).
				Bold(true)

	selectorBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(1).
				BorderForeground(theme.MauveColor)
)

func NewProjectSelector() *ProjectSelectorComponent {
	return &ProjectSelectorComponent{
		mode:            SelectingProject,
		selectedProject: 0,
		selectedTask:    -1,
		selectedTags:    make(map[int]bool),
	}
}

func (c *ProjectSelectorComponent) SetProjects(projects []api.Project) {
	c.projects = projects
	c.selectedProject = 0
}

func (c *ProjectSelectorComponent) SetTasks(tasks []api.Task) {
	c.tasks = tasks
	c.selectedTask = 0
	c.mode = SelectingTask
}

func (c *ProjectSelectorComponent) SetTags(tags []api.Tag) {
	c.tags = tags
	c.selectedTags = make(map[int]bool)
	c.currentTagCursor = 0
}

func (c *ProjectSelectorComponent) SetTagsForEditing(currentTagIDs []string, tags []api.Tag) {
	c.tags = tags
	c.selectedTags = make(map[int]bool)
	c.currentTagCursor = 0
	c.mode = SelectingTags

	for i, tag := range c.tags {
		for _, currentID := range currentTagIDs {
			if tag.ID == currentID {
				c.selectedTags[i] = true
				break
			}
		}
	}
}

func (c *ProjectSelectorComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

func (c *ProjectSelectorComponent) MoveUp() {
	if c.mode == SelectingProject {
		if c.selectedProject > 0 {
			c.selectedProject--
		}
	} else if c.mode == SelectingTask {
		if c.selectedTask > 0 {
			c.selectedTask--
		}
	} else if c.mode == SelectingTags {
		if c.currentTagCursor > 0 {
			c.currentTagCursor--
		}
	}
}

func (c *ProjectSelectorComponent) MoveDown() {
	if c.mode == SelectingProject {
		if c.selectedProject < len(c.projects)-1 {
			c.selectedProject++
		}
	} else if c.mode == SelectingTask {
		if c.selectedTask < len(c.tasks)-1 {
			c.selectedTask++
		}
	} else if c.mode == SelectingTags {
		if c.currentTagCursor < len(c.tags)-1 {
			c.currentTagCursor++
		}
	}
}

func (c *ProjectSelectorComponent) Back() bool {
	if c.mode == SelectingTags {
		c.mode = EnteringDescription
		c.selectedTags = make(map[int]bool)
		c.currentTagCursor = 0
		return false
	}
	if c.mode == EnteringDescription {
		if len(c.tasks) > 0 {
			c.mode = SelectingTask
		} else {
			c.mode = SelectingProject
		}
		c.description = ""
		return false
	}
	if c.mode == SelectingTask {
		c.mode = SelectingProject
		c.selectedTask = -1
		c.tasks = nil
		return false
	}
	return true
}

func (c *ProjectSelectorComponent) GetSelection() (projectID, taskID *string, needsTasks bool) {
	if c.mode == SelectingProject && c.selectedProject >= 0 && c.selectedProject < len(c.projects) {
		pid := c.projects[c.selectedProject].ID
		return &pid, nil, true
	}

	if c.mode == SelectingTask {
		c.mode = EnteringDescription
		return nil, nil, false
	}

	return nil, nil, false
}

func (c *ProjectSelectorComponent) TransitionToTagSelection() {
	if c.mode == EnteringDescription {
		c.mode = SelectingTags
	}
}

func (c *ProjectSelectorComponent) ToggleCurrentTag() {
	if c.mode == SelectingTags && c.currentTagCursor >= 0 && c.currentTagCursor < len(c.tags) {
		c.selectedTags[c.currentTagCursor] = !c.selectedTags[c.currentTagCursor]
	}
}

func (c *ProjectSelectorComponent) ConfirmTags() (projectID, taskID, description *string, tagIDs []string) {
	if c.mode != SelectingTags {
		return nil, nil, nil, nil
	}

	pid := c.projects[c.selectedProject].ID
	var tid *string
	if c.selectedTask >= 0 && c.selectedTask < len(c.tasks) {
		t := c.tasks[c.selectedTask].ID
		tid = &t
	}

	desc := c.description

	var selectedTagIDs []string
	for idx, selected := range c.selectedTags {
		if selected && idx < len(c.tags) {
			selectedTagIDs = append(selectedTagIDs, c.tags[idx].ID)
		}
	}

	return &pid, tid, &desc, selectedTagIDs
}

func (c *ProjectSelectorComponent) GetSelectedTagIDs() []string {
	var selectedTagIDs []string
	for idx, selected := range c.selectedTags {
		if selected && idx < len(c.tags) {
			selectedTagIDs = append(selectedTagIDs, c.tags[idx].ID)
		}
	}
	return selectedTagIDs
}

func (c *ProjectSelectorComponent) ConfirmDescription() (projectID, taskID, description *string) {
	if c.mode != EnteringDescription {
		return nil, nil, nil
	}

	pid := c.projects[c.selectedProject].ID
	var tid *string
	if c.selectedTask >= 0 && c.selectedTask < len(c.tasks) {
		t := c.tasks[c.selectedTask].ID
		tid = &t
	}

	desc := c.description
	return &pid, tid, &desc
}

func (c *ProjectSelectorComponent) AddChar(char rune) {
	if c.mode == EnteringDescription {
		c.description += string(char)
	}
}

func (c *ProjectSelectorComponent) DeleteChar() {
	if c.mode == EnteringDescription && len(c.description) > 0 {
		c.description = c.description[:len(c.description)-1]
	}
}

func (c *ProjectSelectorComponent) GetMode() SelectorMode {
	return c.mode
}

func (c *ProjectSelectorComponent) Reset() {
	c.mode = SelectingProject
	c.selectedProject = 0
	c.selectedTask = -1
	c.selectedTags = make(map[int]bool)
	c.currentTagCursor = 0
	c.tasks = nil
	c.description = ""
}

func (c *ProjectSelectorComponent) GetSelectedProjectID() *string {
	if c.selectedProject >= 0 && c.selectedProject < len(c.projects) {
		pid := c.projects[c.selectedProject].ID
		return &pid
	}
	return nil
}

func (c *ProjectSelectorComponent) View() string {
	if c.mode == SelectingProject {
		return c.renderProjectList()
	} else if c.mode == SelectingTask {
		return c.renderTaskList()
	} else if c.mode == SelectingTags {
		return c.renderTagList()
	}
	return c.renderDescriptionInput()
}

func (c *ProjectSelectorComponent) renderProjectList() string {
	title := selectorTitleStyle.Render("Select Project")
	content := title + "\n\n"

	visibleStart := 0
	visibleEnd := len(c.projects)
	maxVisible := 10

	if len(c.projects) > maxVisible {
		if c.selectedProject > maxVisible/2 {
			visibleStart = c.selectedProject - maxVisible/2
		}
		visibleEnd = visibleStart + maxVisible
		if visibleEnd > len(c.projects) {
			visibleEnd = len(c.projects)
			visibleStart = visibleEnd - maxVisible
			if visibleStart < 0 {
				visibleStart = 0
			}
		}
	}

	for i := visibleStart; i < visibleEnd; i++ {
		project := c.projects[i]
		line := project.Name

		if i == c.selectedProject {
			content += selectorSelectedStyle.Render("▶ "+line) + "\n"
		} else {
			content += selectorItemStyle.Render(line) + "\n"
		}
	}

	content += "\n" + lipgloss.NewStyle().Foreground(theme.Subtext0Color).Render("↑/↓: navigate | enter: select | esc: cancel")

	return selectorBoxStyle.Width(c.width - 4).Render(content)
}

func (c *ProjectSelectorComponent) renderTaskList() string {
	title := selectorTitleStyle.Render("Select Task")
	content := title + "\n\n"

	if len(c.tasks) == 0 {
		content += lipgloss.NewStyle().Foreground(theme.Subtext0Color).Render("No tasks available for this project") + "\n\n"
		content += lipgloss.NewStyle().Foreground(theme.Subtext0Color).Render("enter: continue without task | esc: back")
		return selectorBoxStyle.Width(c.width - 4).Render(content)
	}

	for i, task := range c.tasks {
		line := task.Name

		if i == c.selectedTask {
			content += selectorSelectedStyle.Render("▶ "+line) + "\n"
		} else {
			content += selectorItemStyle.Render(line) + "\n"
		}
	}

	content += "\n" + lipgloss.NewStyle().Foreground(theme.Subtext0Color).Render("↑/↓: navigate | enter: select | esc: back")

	return selectorBoxStyle.Width(c.width - 4).Render(content)
}

func (c *ProjectSelectorComponent) renderDescriptionInput() string {
	title := selectorTitleStyle.Render("Enter Description")
	content := title + "\n\n"

	inputStyle := lipgloss.NewStyle().
		Foreground(theme.GreenColor).
		Bold(true)

	placeholderStyle := lipgloss.NewStyle().
		Foreground(theme.Subtext0Color).
		Italic(true)

	if c.description == "" {
		content += placeholderStyle.Render("Type a description for your timer...") + "\n"
	} else {
		content += inputStyle.Render(c.description) + "█" + "\n"
	}

	content += "\n"

	projectName := "Unknown"
	if c.selectedProject >= 0 && c.selectedProject < len(c.projects) {
		projectName = c.projects[c.selectedProject].Name
	}

	summaryStyle := lipgloss.NewStyle().Foreground(theme.Subtext0Color)
	content += summaryStyle.Render("Project: "+projectName) + "\n"

	if c.selectedTask >= 0 && c.selectedTask < len(c.tasks) {
		taskName := c.tasks[c.selectedTask].Name
		content += summaryStyle.Render("Task: "+taskName) + "\n"
	}

	content += "\n" + lipgloss.NewStyle().Foreground(theme.Subtext0Color).Render("enter: start timer | esc: back")

	return selectorBoxStyle.Width(c.width - 4).Render(content)
}

func (c *ProjectSelectorComponent) renderTagList() string {
	title := selectorTitleStyle.Render("Select Tags (optional)")
	content := title + "\n\n"

	if len(c.tags) == 0 {
		content += lipgloss.NewStyle().Foreground(theme.Subtext0Color).Render("No tags available") + "\n\n"
		content += lipgloss.NewStyle().Foreground(theme.Subtext0Color).Render("enter: continue without tags | esc: back")
		return selectorBoxStyle.Width(c.width - 4).Render(content)
	}

	visibleStart := 0
	visibleEnd := len(c.tags)
	maxVisible := 10

	if len(c.tags) > maxVisible {
		if c.currentTagCursor > maxVisible/2 {
			visibleStart = c.currentTagCursor - maxVisible/2
		}
		visibleEnd = visibleStart + maxVisible
		if visibleEnd > len(c.tags) {
			visibleEnd = len(c.tags)
			visibleStart = visibleEnd - maxVisible
			if visibleStart < 0 {
				visibleStart = 0
			}
		}
	}

	for i := visibleStart; i < visibleEnd; i++ {
		tag := c.tags[i]
		checkbox := "[ ]"
		if c.selectedTags[i] {
			checkbox = "[✓]"
		}

		line := checkbox + " " + tag.Name

		if i == c.currentTagCursor {
			content += selectorSelectedStyle.Render("▶ "+line) + "\n"
		} else {
			content += selectorItemStyle.Render(line) + "\n"
		}
	}

	content += "\n" + lipgloss.NewStyle().Foreground(theme.Subtext0Color).Render("↑/↓: navigate | space: toggle | enter: confirm | esc: back")

	return selectorBoxStyle.Width(c.width - 4).Render(content)
}
