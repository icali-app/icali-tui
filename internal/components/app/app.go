package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/icali-app/icali-tui/internal/style"
)

type Model struct {
	Grid tea.Model
	Preview tea.Model
	EnablePreview bool
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Grid.Init(), m.Preview.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
	width, _ := style.TerminalSize()
	centerGridView := lipgloss.PlaceHorizontal(width, lipgloss.Center, m.Grid.View())

	centerPreviewView := lipgloss.PlaceHorizontal(width, lipgloss.Center, m.Preview.View())

	if m.EnablePreview {
		return lipgloss.JoinVertical(lipgloss.Top, centerGridView, centerPreviewView)
	}
	return centerGridView
}
