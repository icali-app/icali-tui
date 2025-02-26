package grid

import (
	"fmt"
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t":
			c.createTodo()
			return c, c.icsUpdated
		}
	}
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
	t := icshelper.FindTodosForDay(c.info.Calendar, c.info.Day)

	summaryStyle := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(25).
		Height(3)
	
	var summary string
	if len(e) > 0 {
		summary = fmt.Sprintf("Events: %d", len(e))
	} 

	if len(t) > 0 {
		summary = lipgloss.JoinVertical(lipgloss.Top, summary, fmt.Sprintf("Todos: %d", len(t)))
	}

	content += lipgloss.Place(
		width, 
		height,
		lipgloss.Right,
		lipgloss.Bottom,
		summaryStyle.Render(summary),
	)


	return cellStyle.Render(content)
}

// TODO: maybe just make info public
func (d *DayOfMonthCell) Info() DayOfMonthCellInfo {
	return d.info
}


func (c *DayOfMonthCell) createTodo() {
	todo := c.info.Calendar.AddTodo(icshelper.NewId())
	todo.SetStartAt(c.info.Day)
	todo.SetEndAt(c.info.Day)
	todo.SetSummary("new todo of the day")
	todo.SetDescription("some important description")
}

// TODO: This should really be two separate message
// 		 One IcsUpdatedMsg and another CellUpdatedMsg
//		 Since a cell very well might update without changing the ics and vice versa
type IcsUpdatedMsg struct {
	Cell tea.Model
	Calendar *ics.Calendar
}

func (c *DayOfMonthCell) icsUpdated() tea.Msg {
	return IcsUpdatedMsg{Cell: c, Calendar: c.info.Calendar}
}
