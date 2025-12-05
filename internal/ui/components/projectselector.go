package components

import (
	"github.com/charmbracelet/lipgloss"
	"main/internal/api"
)

type SelectorMode int

const (
	SelectingProject SelectorMode = iota
	SelectingTask
	EnteringDescription
)

type ProjectSelectorComponent struct {
	projects        []api.Project
	tasks           []api.Task
	selectedProject int
	selectedTask    int
	mode            SelectorMode
	filterInput     string
	description     string
	width           int
	height          int
}

var (
	selectorTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7C3AED")).
				Padding(0, 1)

	selectorItemStyle = lipgloss.NewStyle().
				PaddingLeft(2)

	selectorSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color("#7C3AED")).
				Bold(true)

	selectorBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(1).
				BorderForeground(lipgloss.Color("#7C3AED"))
)

func NewProjectSelector() *ProjectSelectorComponent {
	return &ProjectSelectorComponent{
		mode:            SelectingProject,
		selectedProject: 0,
		selectedTask:    -1,
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

func (c *ProjectSelectorComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

func (c *ProjectSelectorComponent) MoveUp() {
	if c.mode == SelectingProject {
		if c.selectedProject > 0 {
			c.selectedProject--
		}
	} else {
		if c.selectedTask > 0 {
			c.selectedTask--
		}
	}
}

func (c *ProjectSelectorComponent) MoveDown() {
	if c.mode == SelectingProject {
		if c.selectedProject < len(c.projects)-1 {
			c.selectedProject++
		}
	} else {
		if c.selectedTask < len(c.tasks)-1 {
			c.selectedTask++
		}
	}
}

func (c *ProjectSelectorComponent) Back() bool {
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

	content += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("↑/↓: navigate | enter: select | esc: cancel")

	return selectorBoxStyle.Width(c.width - 4).Render(content)
}

func (c *ProjectSelectorComponent) renderTaskList() string {
	title := selectorTitleStyle.Render("Select Task")
	content := title + "\n\n"

	if len(c.tasks) == 0 {
		content += lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("No tasks available for this project") + "\n\n"
		content += lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("enter: continue without task | esc: back")
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

	content += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("↑/↓: navigate | enter: select | esc: back")

	return selectorBoxStyle.Width(c.width - 4).Render(content)
}

func (c *ProjectSelectorComponent) renderDescriptionInput() string {
	title := selectorTitleStyle.Render("Enter Description")
	content := title + "\n\n"

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Bold(true)

	placeholderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
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

	summaryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	content += summaryStyle.Render("Project: "+projectName) + "\n"

	if c.selectedTask >= 0 && c.selectedTask < len(c.tasks) {
		taskName := c.tasks[c.selectedTask].Name
		content += summaryStyle.Render("Task: "+taskName) + "\n"
	}

	content += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("enter: start timer | esc: back")

	return selectorBoxStyle.Width(c.width - 4).Render(content)
}
