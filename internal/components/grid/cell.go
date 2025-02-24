package grid

import (
	"time"

	ics "github.com/arran4/golang-ical"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	icshelper "github.com/icali-app/icali-tui/internal/ics-helper"
)

// DayOfMonthCell represents a single cell in the grid.
type DayOfMonthCell struct {
	info DayOfMonthCellInfo
}

type DayOfMonthCellInfo struct {
	Day      time.Time
	Calendar *ics.Calendar
}

var (
	width  int = 30
	height int = 5
)

// NewDayOfMonthCell creates a new cell with the provided content.
func NewDayOfMonthCell(info DayOfMonthCellInfo) *DayOfMonthCell {
	return &DayOfMonthCell{
		info: info,
	}
}

// Init implements the tea.Model interface.
func (c *DayOfMonthCell) Init() tea.Cmd {
	// No initialization required for now.
	return nil
}

// Update implements the tea.Model interface.
// Currently, it just returns the component unmodified.
func (c *DayOfMonthCell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

// View implements the tea.Model interface.
// It uses lipgloss to style the cell.
func (c *DayOfMonthCell) View() string {
	// Define a lipgloss style for the cell.

	var content string
	cellStyle := lipgloss.NewStyle().
		Width(width).
		Height(height)

	content += c.info.Day.Format(time.DateOnly)
	e := icshelper.FindEventsForDay(c.info.Calendar, c.info.Day)

	summaryStyle := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(25).
		Height(3)
	
	if len(e) > 0 {
		p := e[0].GetProperty(ics.ComponentPropertySummary)
		content += lipgloss.Place(width, height, lipgloss.Right, lipgloss.Bottom, summaryStyle.Render(p.Value))
	} else {
		content += lipgloss.Place(width, height, lipgloss.Right, lipgloss.Bottom, summaryStyle.Render())
	}

	return cellStyle.Render(content)
}

// TODO: maybe just make info public
func (d *DayOfMonthCell) Info() DayOfMonthCellInfo {
	return d.info
}


// LOGGING
