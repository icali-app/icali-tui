package preview

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/icali-app/icali-tui/internal/components/grid"
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
			m.content = cell.Info().Day.Format(time.DateOnly)
		default:
			panic(fmt.Sprintf("unsupported cell-preview for cell-type: %T", cell))
		}
	}
	return m, nil
}

func (m PreviewComponent) View() string {
	return m.content
}
