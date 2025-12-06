package components

import (
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/lipgloss"
	"main/internal/domain"
	"main/internal/ui/theme"
)

type ReportType int

const (
	DailyReport ReportType = iota
	WeeklyReport
)

type ReportsComponent struct {
	reportType    ReportType
	dailyReport   *domain.DailySummary
	weeklyReport  *domain.WeeklySummary
	selectedDate  time.Time
	width         int
	height        int
}

var (
	reportHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(theme.LavenderColor).
				MarginBottom(1)

	reportTotalStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(theme.GreenColor).
				MarginTop(1)

	reportProjectStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(theme.SapphireColor)

	reportTaskStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			Foreground(theme.Subtext0Color)

	reportBarStyle = lipgloss.NewStyle().
			Foreground(theme.MauveColor)
)

func NewReportsComponent() *ReportsComponent {
	return &ReportsComponent{
		reportType:   DailyReport,
		selectedDate: time.Now(),
	}
}

func (c *ReportsComponent) SetReportType(reportType ReportType) {
	c.reportType = reportType
}

func (c *ReportsComponent) GetReportType() ReportType {
	return c.reportType
}

func (c *ReportsComponent) ToggleReportType() {
	if c.reportType == DailyReport {
		c.reportType = WeeklyReport
	} else {
		c.reportType = DailyReport
	}
}

func (c *ReportsComponent) SetDailyReport(report *domain.DailySummary) {
	c.dailyReport = report
}

func (c *ReportsComponent) SetWeeklyReport(report *domain.WeeklySummary) {
	c.weeklyReport = report
}

func (c *ReportsComponent) SetSelectedDate(date time.Time) {
	c.selectedDate = date
}

func (c *ReportsComponent) GetSelectedDate() time.Time {
	return c.selectedDate
}

func (c *ReportsComponent) NextDate() {
	if c.reportType == DailyReport {
		c.selectedDate = c.selectedDate.AddDate(0, 0, 1)
	} else {
		c.selectedDate = c.selectedDate.AddDate(0, 0, 7)
	}
}

func (c *ReportsComponent) PrevDate() {
	if c.reportType == DailyReport {
		c.selectedDate = c.selectedDate.AddDate(0, 0, -1)
	} else {
		c.selectedDate = c.selectedDate.AddDate(0, 0, -7)
	}
}

func (c *ReportsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

func (c *ReportsComponent) View() string {
	if c.reportType == DailyReport {
		return c.renderDailyReport()
	}
	return c.renderWeeklyReport()
}

func (c *ReportsComponent) renderDailyReport() string {
	if c.dailyReport == nil {
		return lipgloss.NewStyle().Foreground(theme.Subtext0Color).
			Render("Loading daily report...")
	}

	dateStr := c.dailyReport.Date.Format("Monday, January 2, 2006")
	content := reportHeaderStyle.Render(dateStr) + "\n\n"

	if c.dailyReport.TotalDuration == 0 {
		content += lipgloss.NewStyle().Foreground(theme.Subtext0Color).
			Italic(true).Render("No time tracked on this day")
		return content
	}

	projects := make([]string, 0, len(c.dailyReport.ByProject))
	for projectID := range c.dailyReport.ByProject {
		projects = append(projects, projectID)
	}
	sort.Slice(projects, func(i, j int) bool {
		return c.dailyReport.ByProject[projects[i]].TotalDuration >
			c.dailyReport.ByProject[projects[j]].TotalDuration
	})

	for _, projectID := range projects {
		projectSummary := c.dailyReport.ByProject[projectID]
		projectLine := fmt.Sprintf("%s - %s",
			projectSummary.ProjectName,
			domain.FormatDuration(projectSummary.TotalDuration))
		content += reportProjectStyle.Render(projectLine) + "\n"

		tasks := make([]string, 0, len(projectSummary.ByTask))
		for taskID := range projectSummary.ByTask {
			tasks = append(tasks, taskID)
		}
		sort.Slice(tasks, func(i, j int) bool {
			return projectSummary.ByTask[tasks[i]].Duration >
				projectSummary.ByTask[tasks[j]].Duration
		})

		for _, taskID := range tasks {
			taskSummary := projectSummary.ByTask[taskID]
			taskLine := fmt.Sprintf("  • %s - %s",
				taskSummary.TaskName,
				domain.FormatDuration(taskSummary.Duration))
			content += reportTaskStyle.Render(taskLine) + "\n"
		}

		content += "\n"
	}

	totalLine := fmt.Sprintf("Total: %s", domain.FormatDuration(c.dailyReport.TotalDuration))
	content += reportTotalStyle.Render(totalLine)

	helpText := "\n\n" + lipgloss.NewStyle().Foreground(theme.Subtext0Color).
		Render("←/→: prev/next day | t: toggle report type")

	return content + helpText
}

func (c *ReportsComponent) renderWeeklyReport() string {
	if c.weeklyReport == nil {
		return lipgloss.NewStyle().Foreground(theme.Subtext0Color).
			Render("Loading weekly report...")
	}

	weekStr := fmt.Sprintf("Week of %s - %s",
		c.weeklyReport.StartDate.Format("Jan 2"),
		c.weeklyReport.EndDate.AddDate(0, 0, -1).Format("Jan 2, 2006"))
	content := reportHeaderStyle.Render(weekStr) + "\n\n"

	if c.weeklyReport.TotalDuration == 0 {
		content += lipgloss.NewStyle().Foreground(theme.Subtext0Color).
			Italic(true).Render("No time tracked this week")
		return content
	}

	content += lipgloss.NewStyle().Bold(true).Render("Daily Breakdown:") + "\n"
	for i := 0; i < 7; i++ {
		date := c.weeklyReport.StartDate.AddDate(0, 0, i)
		dateKey := date.Format("2006-01-02")
		duration := c.weeklyReport.ByDay[dateKey]

		dayName := date.Format("Mon")
		dayStr := date.Format("Jan 2")

		if duration > 0 {
			bar := c.createBar(duration, c.weeklyReport.TotalDuration, 20)
			line := fmt.Sprintf("  %s %s: %s %s",
				dayName, dayStr,
				domain.FormatDuration(duration),
				bar)
			content += line + "\n"
		} else {
			line := fmt.Sprintf("  %s %s: -", dayName, dayStr)
			content += lipgloss.NewStyle().Foreground(theme.Subtext0Color).Render(line) + "\n"
		}
	}

	content += "\n" + lipgloss.NewStyle().Bold(true).Render("By Project:") + "\n"

	projects := make([]string, 0, len(c.weeklyReport.ByProject))
	for projectID := range c.weeklyReport.ByProject {
		projects = append(projects, projectID)
	}
	sort.Slice(projects, func(i, j int) bool {
		return c.weeklyReport.ByProject[projects[i]].TotalDuration >
			c.weeklyReport.ByProject[projects[j]].TotalDuration
	})

	for _, projectID := range projects {
		projectSummary := c.weeklyReport.ByProject[projectID]
		projectLine := fmt.Sprintf("  %s - %s",
			projectSummary.ProjectName,
			domain.FormatDuration(projectSummary.TotalDuration))
		content += reportProjectStyle.Render(projectLine) + "\n"

		tasks := make([]string, 0, len(projectSummary.ByTask))
		for taskID := range projectSummary.ByTask {
			tasks = append(tasks, taskID)
		}
		sort.Slice(tasks, func(i, j int) bool {
			return projectSummary.ByTask[tasks[i]].Duration >
				projectSummary.ByTask[tasks[j]].Duration
		})

		for _, taskID := range tasks {
			taskSummary := projectSummary.ByTask[taskID]
			taskLine := fmt.Sprintf("    • %s - %s",
				taskSummary.TaskName,
				domain.FormatDuration(taskSummary.Duration))
			content += reportTaskStyle.Render(taskLine) + "\n"
		}
	}

	totalLine := fmt.Sprintf("\nTotal: %s", domain.FormatDuration(c.weeklyReport.TotalDuration))
	content += reportTotalStyle.Render(totalLine)

	helpText := "\n\n" + lipgloss.NewStyle().Foreground(theme.Subtext0Color).
		Render("←/→: prev/next week | t: toggle report type")

	return content + helpText
}

func (c *ReportsComponent) createBar(value, max time.Duration, maxWidth int) string {
	if max == 0 {
		return ""
	}

	percentage := float64(value) / float64(max)
	barWidth := int(percentage * float64(maxWidth))
	if barWidth < 1 && value > 0 {
		barWidth = 1
	}

	bar := ""
	for i := 0; i < barWidth; i++ {
		bar += "█"
	}

	return reportBarStyle.Render(bar)
}
