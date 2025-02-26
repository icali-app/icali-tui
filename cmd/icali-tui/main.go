package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"time"

	ics "github.com/arran4/golang-ical"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/icali-app/icali-tui/internal/components/app"
	"github.com/icali-app/icali-tui/internal/components/grid"
	"github.com/icali-app/icali-tui/internal/components/preview"
	"github.com/icali-app/icali-tui/internal/components/toast"
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


type TestStorage struct {
	path string
}

func (s *TestStorage) Upload(data []byte) error {
	return os.WriteFile(s.path, data, 0644)
}

func (s *TestStorage) Download() ([]byte, error) {
	file, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}


func main() {
	_ = config.Get()

	p := os.Args[1]
	s := TestStorage{path: p}
	data, err := s.Download()
	if err != nil {
		fmt.Println("Failed to open calendar:", err)
		os.Exit(1)
	}

	calendar, err = ics.ParseCalendar(bytes.NewReader(data))
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
		Storage: &s,
	}

	appWithToast := toast.WithToast(app)

	if _, err := tea.NewProgram(appWithToast, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
