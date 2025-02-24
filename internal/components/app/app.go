package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Grid tea.Model
	Preview tea.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.Grid, cmd = m.Grid.Update(msg)

	m.Preview, cmd = m.Preview.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.Grid.View() + "\n" + m.Preview.View()
}
