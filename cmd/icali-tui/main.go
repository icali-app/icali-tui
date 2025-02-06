package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/icali-app/icali-tui/internal/components/app"
)

func main() {
	app := app.Model{}
	if _, err := tea.NewProgram(app).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
