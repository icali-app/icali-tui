package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Grid tea.Model
	Preview tea.Model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Grid.Init(), m.Preview.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmdGrid, cmdPreview tea.Cmd
	m.Grid, cmdGrid = m.Grid.Update(msg)
	m.Preview, cmdPreview = m.Preview.Update(msg)
	return m, tea.Batch(cmdGrid, cmdPreview)
}

func (m Model) View() string {
	return m.Grid.View() + "\n" + m.Preview.View()
}
