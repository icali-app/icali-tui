package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/icali-app/icali-tui/internal/config"
	"golang.org/x/term"
	"os"
	"sync"
)

var (
	conf  = config.Get()
	once  sync.Once
	style Style
)

type Style struct {
	Base             lipgloss.Style
	WithSelectedText lipgloss.Style
	WithBorder       lipgloss.Style
	WithLink         lipgloss.Style
	WithSummary      lipgloss.Style
}

func Get() Style {
	once.Do(func() {
		style = Style{
			Base: lipgloss.NewStyle().
				Background(lipgloss.Color(conf.Style.Normal.Background)).
				Foreground(lipgloss.Color(conf.Style.Normal.Text)),

			WithSelectedText: lipgloss.NewStyle().
				Inherit(style.Base).
				Foreground(lipgloss.Color(conf.Style.Normal.Selection)).
				Bold(true),

			WithBorder: lipgloss.NewStyle().
				Inherit(style.Base).
				Border(lipgloss.RoundedBorder()).
				BorderBackground(lipgloss.Color(conf.Style.Normal.Background)).
				BorderForeground(lipgloss.Color(conf.Style.Normal.Border)),

			WithLink: lipgloss.NewStyle().
				Inherit(style.Base).
				Foreground(lipgloss.Color(conf.Style.Normal.Link)),

			WithSummary: lipgloss.NewStyle().
				Inherit(style.Base).
				Bold(true),
		}
	})

	return style
}

func TerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic("Could not get terminal size")
	}

	return width, height
}
