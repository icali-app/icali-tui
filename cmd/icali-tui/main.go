package main

import (
	"fmt"
	"os"

	"time"

	ics "github.com/arran4/golang-ical"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/icali-app/icali-tui/internal/components/app"
	"github.com/icali-app/icali-tui/internal/components/grid"
	"github.com/icali-app/icali-tui/internal/components/preview"
	"github.com/icali-app/icali-tui/internal/config"
)

var (
	calendar *ics.Calendar
)

func formatDate(t time.Time) string {
	return t.Format(time.DateOnly)
}

func createExampleCell(row, col, cursor int) *grid.DayOfMonthCell {
	now := time.Date(2024, time.October, 1, 0, 0, 0, 0, time.UTC)
	cellDate := now.AddDate(0, 0, cursor)
	info := grid.DayOfMonthCellInfo{
		Day:      cellDate,
		Calendar: calendar,
	}
	return grid.NewDayOfMonthCell(info)
}

func main() {
	_ = config.Get()

	p := os.Args[1]
	file, err := os.Open(p)
	if err != nil {
		fmt.Println("Failed to open calendar:", err)
		os.Exit(1)
	}

	calendar, err = ics.ParseCalendar(file)
	if err != nil {
		fmt.Println("Failed to parse calendar:", err)
		os.Exit(1)
	}

	grid := grid.NewGridComponentWithCellFunc(3, 4, func(row, col, cursor int) tea.Model {
		return createExampleCell(row, col, cursor)
	})

	pr := preview.NewPreview()

	app := app.Model{
		Grid: grid,
		Preview: pr,
		EnablePreview: true,
	}
	if _, err := tea.NewProgram(app, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
