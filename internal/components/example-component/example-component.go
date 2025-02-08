package examplecomponent

import (
	"fmt"
	"os"

	ics "github.com/arran4/golang-ical"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cal          *ics.Calendar
	currentEntry int
}

func NewFromCalPathStr(path string) (Model, error) {
	f, err := os.Open(path)
	if err != nil {
		return Model{}, err
	}

	cal, err := ics.ParseCalendar(f)
	if err != nil {
		return Model{}, err
	}

	err = f.Close()
	if err != nil {
		return Model{}, err
	}

	fmt.Println("cal created")

	return Model{cal, 0}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.currentEntry > 0 {
				m.currentEntry--
			}
		case "down", "j":
			if m.currentEntry < 4-1 {
				m.currentEntry++
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	b := []string{"hello", "du", "lustiger", "spasst"}
	a := fmt.Sprintf("%d: %+v", m.currentEntry, b[m.currentEntry])
	return a
}
