package preview

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/icali-app/icali-tui/internal/components/grid"
	icshelper "github.com/icali-app/icali-tui/internal/ics-helper"
	"github.com/icali-app/icali-tui/internal/style"
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
		p := e.GetProperty(ics.ComponentPropertySummary)
		styledP := styl.WithBorder.Render(p.Value)
		eventsStr = lipgloss.JoinHorizontal(lipgloss.Top, eventsStr, styledP)
	}

	return lipgloss.JoinVertical(lipgloss.Top, date, eventsStr)
}

func (m PreviewComponent) View() string {
	return styl.WithBorder.Render(m.content)
}
