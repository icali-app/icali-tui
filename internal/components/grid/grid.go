package grid

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GridMode int

const (
	GlobalMode GridMode = 0
	InsertMode GridMode = 1
)

// GridComponent represents a grid of CellComponents.
type GridComponent struct {
	rows int
	cols int
	// The cursor is 1-D although a grid is 2-D
	// Therefore, the cursor position = i * rows + col
	// Why? Because this allows for easy restructoring of the grid if needed (e.g. remove a column and add rows)
	cursor int
	cells  [][]*CellComponent // Todo flatten this array too
	mode   GridMode
}

// NewGridComponent creates a new grid with the given number of rows and columns.
func NewGridComponent(rows, cols int) *GridComponent {
	grid := &GridComponent{
		rows:  rows,
		cols:  cols,
		cells: make([][]*CellComponent, rows),
	}

	for i := 0; i < rows; i++ {
		grid.cells[i] = make([]*CellComponent, cols)
		for j := 0; j < cols; j++ {
			// For demonstration, each cell shows its coordinate.
			content := fmt.Sprintf("(%d,%d)", i, j)
			grid.cells[i][j] = NewCellComponent(content)
		}
	}

	return grid
}

// Init initializes all cell components within the grid.
func (g *GridComponent) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, row := range g.cells {
		for _, cell := range row {
			cmds = append(cmds, cell.Init())
		}
	}
	return tea.Batch(cmds...)
}

// Update propagates incoming messages to each cell in the grid.
func (g *GridComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch g.mode {
	case GlobalMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				g.moveUp()
			case "down", "j":
				g.moveDown()
			case "left", "h":
				g.moveLeft()
			case "right", "l":
				g.moveRight()
			case "i":
				g.mode = InsertMode
			}
		}
	case InsertMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				g.mode = GlobalMode
			default:
				cell := g.currentCell()
				_, cmd := cell.Update(msg)
				return g, cmd
			}
		}
	}
	// var cmds []tea.Cmd
	// for i, row := range g.cells {
	// 	for j, cell := range row {
	// 		updated, cmd := cell.Update(msg)
	// 		g.cells[i][j] = updated.(*CellComponent)
	// 		if cmd != nil {
	// 			cmds = append(cmds, cmd)
	// 		}
	// 	}
	// }
	// return g, tea.Batch(cmds...)

	return g, nil
}

// View renders the grid by joining the cells using lipgloss.
func (g *GridComponent) View() string {
	var rows []string
	for ridx, row := range g.cells {
		var cellViews []string
		for cidx, cell := range row {
			content := cell.View()

			var style lipgloss.Style
			if g.isCursorAt(ridx, cidx) {
				style = lipgloss.NewStyle().
					Border(lipgloss.ThickBorder()).
					BorderForeground(lipgloss.Color("#FF0000")).
					Padding(3)
			} else {
				style = lipgloss.NewStyle().
					Border(lipgloss.NormalBorder()).
					Padding(3)
			}

			cellViews = append(cellViews, style.Render(content))
		}
		// Join all cells in a row horizontally.
		rowStr := lipgloss.JoinHorizontal(lipgloss.Top, cellViews...)
		rows = append(rows, rowStr)
	}
	// Join all rows vertically to form the grid.
	gridView := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return gridView
}

func (g *GridComponent) moveUp() {
	for i := 0; i < g.cols; i++ {
		g.moveLeft()
	}
}

func (g *GridComponent) moveDown() {
	for i := 0; i < g.cols; i++ {
		g.moveRight()
	}
}

func (g *GridComponent) moveLeft() {
	if g.cursor > 0 {
		g.cursor--
	}
}

func (g *GridComponent) moveRight() {
	if g.cursor < g.cellCount()-1 {
		g.cursor++
	}
}

func (g *GridComponent) cellCount() int {
	return g.rows * g.cols
}

func (g *GridComponent) isCursorAt(ridx, cidx int) bool {
	return g.cursor == ridx*g.cols+cidx
}

func (g *GridComponent) currentPos() (int, int) {
	cidx := g.cursor % g.cols
	ridx := g.cursor / g.cols
	return ridx, cidx
}

func (g *GridComponent) currentCell() *CellComponent {
	ridx, cidx := g.currentPos()
	return g.cells[ridx][cidx]
}
