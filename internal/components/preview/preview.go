package preview

import (
	tea "github.com/charmbracelet/bubbletea"
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
	return m, nil
}

func (m PreviewComponent) View() string {
	return m.content
}
