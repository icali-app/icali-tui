package grid

import (
	"time"

	ics "github.com/arran4/golang-ical"
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
	// The position of the curor is equal to some row,col with cursor = row * cols + col
	// Where row is in [0, rows) and col is in [0, cols)
	rows   int
	cols   int
	cursor int
	cells  []tea.Model
	mode   GridMode
}

// NewGridComponent creates a new grid with the given number of rows and columns.
func NewGridComponent(rows, cols int) *GridComponent {
	grid := &GridComponent{
		rows:  rows,
		cols:  cols,
		cells: make([]tea.Model, rows*cols),
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			// For demonstration, each cell shows its coordinate.
			info := DayOfMonthCellInfo{
				Day:      time.Now(),
				Calendar: ics.NewCalendar(),
			}
			grid.cells[row*cols+col] = NewDayOfMonthCell(info)
		}
	}

	return grid
}

type CellFunc = func(row, col, cursor int) tea.Model 

func NewGridComponentWithCellFunc(rows, cols int, cellFunc CellFunc) *GridComponent {
	grid := &GridComponent{
		rows:  rows,
		cols:  cols,
		cells: make([]tea.Model, rows*cols),
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			// For demonstration, each cell shows its coordinate.
			grid.cells[row*cols+col] = cellFunc(row, col, row*cols+col)
		}
	}

	return grid
}

// Init initializes all cell components within the grid.
func (g *GridComponent) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, cell := range g.cells {
		cmds = append(cmds, cell.Init())
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
	var(
		normalStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder())

		selectedStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#FF0000"))

		outerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())
	)

	var rows []string
	for row := 0; row < g.rows; row++ {
		var cellViews []string
		for col := 0; col < g.cols; col++ {
			cell := g.cellAt(row, col)
			content := cell.View()

			var style lipgloss.Style
			if g.isCursorAt(row, col) {
				style = selectedStyle
			} else {
				style = normalStyle
			}

			// Disable all border to selectively enable them
			style = style.
				BorderBottom(false).
				BorderTop(false).
				BorderLeft(false).
				BorderRight(false)

			if col != g.cols - 1 {
				style = style.BorderRight(true)
			}

			if row != g.rows - 1 {
				style = style.BorderBottom(true)
			}

			cellViews = append(cellViews, style.Render(content))
		}
		// Join all cells in a row horizontally.
		rowStr := lipgloss.JoinHorizontal(lipgloss.Top, cellViews...)
		rows = append(rows, rowStr)
	}
	// Join all rows vertically to form the grid.
	gridView := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return outerStyle.Render(gridView)
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

func (g *GridComponent) isCursorAt(row, col int) bool {
	return g.cursor == row*g.cols+col
}

func (g *GridComponent) currentPos() (int, int) {
	col := g.cursor % g.cols
	row := g.cursor / g.cols
	return row, col
}

func (g *GridComponent) currentCell() tea.Model {
	row, col := g.currentPos()
	return g.cellAt(row, col)
}

func (g *GridComponent) cellAt(row, col int) tea.Model {
	return g.cells[row*g.cols+col]
}
