package main

import (
	"fmt"
	"os"
	"time"

	ics "github.com/arran4/golang-ical"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/icali-app/icali-tui/internal/components/app"
	"github.com/icali-app/icali-tui/internal/components/grid"
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
	// p := os.Args[1]
	// fmt.Println("Starting app")
	// child, err := examplecomponent.NewFromCalPathStr(p)
	// if err != nil {
	// 	panic(err)
	// }

	_ = config.Get()

	p := os.Args[1]
	file, err := os.Open(p)
	if err != nil {
		panic(err)
	}

	calendar, err = ics.ParseCalendar(file)
	if err != nil {
		panic(err)
	}

	grid := grid.NewGridComponentWithCellFunc(3, 4, func(row, col, cursor int) tea.Model {
		return createExampleCell(row, col, cursor)
	})
	app := app.Model{
		Grid: grid,
	}
	if _, err := tea.NewProgram(app).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
