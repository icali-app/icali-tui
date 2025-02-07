package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/icali-app/icali-tui/internal/components/app"
	examplecomponent "github.com/icali-app/icali-tui/internal/components/example-component"
)

func main() {
	p := os.Args[1]
	fmt.Println("Starting app")
	child, err := examplecomponent.NewFromCalPathStr(p)
	if err != nil {
		panic(err)
	}
	app := app.Model{
		Haha: child,
	}
	if _, err := tea.NewProgram(app).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
