package views

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"main/internal/domain"
	"main/internal/ui/components"
	"main/internal/ui/theme"
)

type ReportsView struct {
	reportsComponent *components.ReportsComponent
	width            int
	height           int
}

func NewReportsView() *ReportsView {
	return &ReportsView{
		reportsComponent: components.NewReportsComponent(),
	}
}

func (v *ReportsView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.reportsComponent.SetSize(width, height)
}

func (v *ReportsView) GetReportType() components.ReportType {
	return v.reportsComponent.GetReportType()
}

func (v *ReportsView) GetSelectedDate() time.Time {
	return v.reportsComponent.GetSelectedDate()
}

func (v *ReportsView) SetDailyReport(report *domain.DailySummary) {
	v.reportsComponent.SetDailyReport(report)
}

func (v *ReportsView) SetWeeklyReport(report *domain.WeeklySummary) {
	v.reportsComponent.SetWeeklyReport(report)
}

func (v *ReportsView) SetTags(tags map[string]string) {
	v.reportsComponent.SetTags(tags)
}

func (v *ReportsView) ToggleReportType() {
	v.reportsComponent.ToggleReportType()
}

func (v *ReportsView) NextDate() {
	v.reportsComponent.NextDate()
}

func (v *ReportsView) PrevDate() {
	v.reportsComponent.PrevDate()
}

func (v *ReportsView) Update(msg tea.Msg) (*ReportsView, tea.Cmd) {
	return v, nil
}

func (v *ReportsView) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.LavenderColor).
		MarginBottom(1)

	reportTypeStr := "Daily"
	if v.reportsComponent.GetReportType() == components.WeeklyReport {
		reportTypeStr = "Weekly"
	}

	content := titleStyle.Render("📊 Reports - "+reportTypeStr) + "\n\n"
	content += v.reportsComponent.View()

	return content
}
