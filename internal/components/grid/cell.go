package grid

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CellComponent represents a single cell in the grid.
type CellComponent struct {
	content string
}

// NewCellComponent creates a new cell with the provided content.
func NewCellComponent(content string) *CellComponent {
	return &CellComponent{
		content: content,
	}
}

// Init implements the tea.Model interface.
func (c *CellComponent) Init() tea.Cmd {
	// No initialization required for now.
	return nil
}

// Update implements the tea.Model interface.
// Currently, it just returns the component unmodified.
func (c *CellComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

// View implements the tea.Model interface.
// It uses lipgloss to style the cell.
func (c *CellComponent) View() string {
	// Define a lipgloss style for the cell.
	cellStyle := lipgloss.NewStyle()
	return cellStyle.Render(c.content)
}
