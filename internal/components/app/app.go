package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/icali-app/icali-tui/internal/components/grid"
	"github.com/icali-app/icali-tui/internal/components/toast"
	"github.com/icali-app/icali-tui/internal/storage"
	"github.com/icali-app/icali-tui/internal/style"
)

type Model struct {
	Grid tea.Model
	Preview tea.Model
	EnablePreview bool
	Storage storage.Storage
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Grid.Init(), m.Preview.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case grid.IcsUpdatedMsg: // TODO: maybe refactor this to its own package
		cal := msg.Calendar
		bytes := []byte(cal.Serialize())
		err := m.Storage.Upload(bytes)
		if err != nil {	
			err = fmt.Errorf("Failed to save calendar: %w", err)
			return m, toast.Error(err.Error())
		}

		return m, toast.Success("Calendar saved.")

	case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
			return m, tea.Quit
		case "p":
			m.EnablePreview = !m.EnablePreview
		}
	}

	var cmdGrid, cmdPreview tea.Cmd
	m.Grid, cmdGrid = m.Grid.Update(msg)
	m.Preview, cmdPreview = m.Preview.Update(msg)
	return m, tea.Batch(cmdGrid, cmdPreview)
}

func (m Model) View() string {
	width, height := style.TerminalSize()
	fullscreenStyle := lipgloss.NewStyle().
		Width(width).
		Height(height)

	centerGridView := lipgloss.PlaceHorizontal(width, lipgloss.Center, m.Grid.View())

	centerPreviewView := lipgloss.PlaceHorizontal(width, lipgloss.Center, m.Preview.View())

	var res string
	if m.EnablePreview {
		res = lipgloss.JoinVertical(lipgloss.Top, centerGridView, centerPreviewView)
	}

	return fullscreenStyle.Render(res)

}
