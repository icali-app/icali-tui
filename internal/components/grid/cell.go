package grid

import (
	"time"

	ics "github.com/arran4/golang-ical"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	e := c.findEventsForThisDay()

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

func (cell *DayOfMonthCell) findEventsForThisDay() []*ics.VEvent {
	res := make([]*ics.VEvent, 0)

	for _, event := range cell.info.Calendar.Events() {
		start, err := event.GetStartAt()

		if err == nil && isSameDay(cell.info.Day, start) {
			res = append(res, event)
			continue
		}

		// TODO: we will ignore everything else for now
		// end, err := event.GetStartAt()
		// if err == nil && isSameDay(cell.info.day, end) {
		// 	res = append(res, event)
		// 	continue
		// }
	}

	return res
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}


// LOGGING

