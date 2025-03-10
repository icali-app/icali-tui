package preview

import (
	"fmt"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/icali-app/icali-tui/internal/components/grid"
	"github.com/icali-app/icali-tui/internal/ellipsis"
	icshelper "github.com/icali-app/icali-tui/internal/ics-helper"
	"github.com/icali-app/icali-tui/internal/style"
	"github.com/icali-app/icali-tui/internal/tiss"
)

var (
	styl = style.Get()
)

type PreviewComponent struct {
	currentCell tea.Model
}

func NewPreview() *PreviewComponent {
	return &PreviewComponent{nil}
}

func (m PreviewComponent) Init() tea.Cmd {
	return nil
}

func (m PreviewComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case grid.SelectedCellMsg:
		cell := msg.Cell
		switch cell := cell.(type) {
		case *grid.DayOfMonthCell:	
			m.currentCell = cell
		default:
			panic(fmt.Sprintf("unsupported cell-preview for cell-type: %T", cell))
		}
	case grid.IcsUpdatedMsg:
		cell := msg.Cell
		if cell == m.currentCell {
			switch cell := cell.(type) {
			case *grid.DayOfMonthCell:	
				m.currentCell = cell
			default:
				panic(fmt.Sprintf("unsupported cell-preview for cell-type: %T", cell))
			}
		}
	}
	return m, nil
}

func formatDayOfMonthCell(cell grid.DayOfMonthCell) string {
	info := cell.Info()
	date := fmt.Sprintf("Date: %s", info.Day.Format(time.DateOnly))
	events := icshelper.FindEventsForDay(info.Calendar, info.Day)

	var eventsStr string	
	for _, e := range events {
		summary := e.GetProperty(ics.ComponentPropertySummary).Value
		summary = styl.WithSummary.Render(summary)

		location := e.GetProperty(ics.ComponentPropertyLocation).Value
		locationText := fmt.Sprintf("Location: %s", ellipsis.WithEllipsis(location, 50))

		locationLink, err := locationLink(location, e)
		if err != nil {
			location = styl.WithInvalidLink.Render(locationText)
		} else {
			locationFmt := formatHyperlink(locationLink, locationText)
			location = styl.WithLink.Render(locationFmt)
		}


		description := e.GetProperty(ics.ComponentPropertyDescription).Value
		description = styl.Base.Render(description)

		styled := lipgloss.JoinVertical(lipgloss.Top, summary, location, description)
		styled = styl.WithBorder.Render(styled)
		eventsStr = lipgloss.JoinHorizontal(lipgloss.Top, eventsStr, styled)
	}

	var todosStr string
	todos := icshelper.FindTodosForDay(info.Calendar, info.Day)
	for _, t := range todos {
		summary := t.GetProperty(ics.ComponentPropertySummary).Value
		summary = styl.WithSummary.Render(summary)

		description := t.GetProperty(ics.ComponentPropertyDescription).Value
		description = styl.Base.Render(description)

		styled := lipgloss.JoinVertical(lipgloss.Top, summary, description)
		styled = styl.WithBorder.Render(styled)
		eventsStr = lipgloss.JoinHorizontal(lipgloss.Top, eventsStr, styled)
	}

	return lipgloss.JoinVertical(lipgloss.Top, date, eventsStr, todosStr)
}

func formatCell(cell tea.Model) string {
	if cell == nil {
		return "nothing selected..."
	}

 	switch cell := cell.(type) {
	case *grid.DayOfMonthCell:	
		return formatDayOfMonthCell(*cell)
	default:
		panic(fmt.Sprintf("unsupported cell-preview for cell-type: %T", cell))
	}
}

func (m PreviewComponent) View() string {
	content := formatCell(m.currentCell)
	return styl.WithBorder.Render(content)
}


func formatHyperlink(link, text string) string {
	return fmt.Sprintf("\x1B]8;;%s\x1B\\%s\x1B]8;;\x1B\\\n", link, text)
}

func locationLink(str string, reference *ics.VEvent) (string, error) {
	if strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://") {
		return str, nil
	}

	// TU Wien calendar check
	uid := reference.GetProperty(ics.ComponentPropertyUniqueId).Value
	if strings.HasSuffix(uid, "tuwien.ac.at") {
		location := reference.GetProperty(ics.ComponentPropertyLocation).Value
		return tiss.GetTissRoom(location)
	}

	return "http://example.com", nil
}
