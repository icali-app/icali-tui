package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/icali-app/icali-tui/internal/components/app"
	"github.com/icali-app/icali-tui/internal/components/grid"
)

func main() {
	// p := os.Args[1]
	// fmt.Println("Starting app")
	// child, err := examplecomponent.NewFromCalPathStr(p)
	// if err != nil {
	// 	panic(err)
	// }

	grid := grid.NewGridComponent(3, 4)
	app := app.Model{
		Grid: grid,
	}
	if _, err := tea.NewProgram(app).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
