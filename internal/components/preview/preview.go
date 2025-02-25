package preview

import (
	"fmt"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/icali-app/icali-tui/internal/components/grid"
	icshelper "github.com/icali-app/icali-tui/internal/ics-helper"
	"github.com/icali-app/icali-tui/internal/style"
	"github.com/icali-app/icali-tui/internal/tiss"
)

var (
	styl = style.Get()
)

type PreviewComponent struct {
	content string
}

func NewPreview(content string) *PreviewComponent {
	return &PreviewComponent{content}
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
			m.content = formatDayOfMonthCell(*cell)
		default:
			panic(fmt.Sprintf("unsupported cell-preview for cell-type: %T", cell))
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
		locationLink := locationLink(location, e)
		locationText := fmt.Sprintf("Location (%s)", location)
		locationFmt := formatHyperlink(locationLink, locationText)
		location = styl.WithLink.Render(locationFmt)

		description := e.GetProperty(ics.ComponentPropertyDescription).Value
		description = styl.Base.Render(description)

		styled := lipgloss.JoinVertical(lipgloss.Top, summary, location, description)
		styled = styl.WithBorder.Render(styled)
		eventsStr = lipgloss.JoinHorizontal(lipgloss.Top, eventsStr, styled)
	}

	return lipgloss.JoinVertical(lipgloss.Top, date, eventsStr)
}

func (m PreviewComponent) View() string {
	return styl.WithBorder.Render(m.content)
}


func formatHyperlink(link, text string) string {
	return fmt.Sprintf("\x1B]8;;%s\x1B\\%s\x1B]8;;\x1B\\\n", link, text)
}

func locationLink(str string, reference *ics.VEvent) string {
	if strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://") {
		return str
	}

	// TU Wien calendar check
	uid := reference.GetProperty(ics.ComponentPropertyUniqueId).Value
	if strings.HasSuffix(uid, "tuwien.ac.at") {
		location := reference.GetProperty(ics.ComponentPropertyLocation).Value
		return tiss.SearchTISSRoom(location)
	}

	return "http://example.com"
}
