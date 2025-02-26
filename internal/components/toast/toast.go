package toast

import (
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/icali-app/icali-tui/internal/config"
	"github.com/icali-app/icali-tui/internal/logger"
	"github.com/icali-app/icali-tui/internal/style"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

var (
	conf                     = config.Get()
	styl                     = style.Get()
	AppToast   *Toast        = &Toast{msg: "bla"}
	defaultTTL time.Duration = time.Second * 5
)

type Toast struct {
	msg   string
	style lipgloss.Style
	ttl   time.Duration
}

func New(msg string) *Toast {
	return &Toast{msg, styl.WithBorder, defaultTTL}
}

func NewSuccess(msg string) *Toast {
	successBorder := lipgloss.NewStyle().
		Inherit(styl.WithBorder).
		BorderForeground(lipgloss.Color(conf.Style.Success.Border))
	return &Toast{msg, successBorder, defaultTTL}
}

func NewError(msg string) *Toast {
	errorBorder := lipgloss.NewStyle().
		Inherit(styl.WithBorder).
		BorderForeground(lipgloss.Color(conf.Style.Error.Border))
	return &Toast{msg, errorBorder, defaultTTL}
}

func NewInfo(msg string) *Toast {
	infoBorder := lipgloss.NewStyle().
		Inherit(styl.WithBorder).
		BorderForeground(lipgloss.Color(conf.Style.Info.Border))
	return &Toast{msg, infoBorder, defaultTTL}
}

type toasted struct {
	model       tea.Model
	activeToasts []tea.Model
}

func (t *toasted) push(toast tea.Model) {
	t.activeToasts = append(t.activeToasts, toast)
}

func (t *toasted) pop(toast tea.Model) {
	t.activeToasts = slices.DeleteFunc(t.activeToasts, func(activeToast tea.Model) bool {
		return activeToast == toast
	})
}

func (t toasted) Init() tea.Cmd {
	return t.model.Init()
}

func (t toasted) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l := logger.Get()
	switch msg := msg.(type) {
	case toastDeadMsg:
		t.pop(msg.toast)
		l.Debug().Int("toast-count", len(t.activeToasts)).Msg("Removed toast")
	case viewToastMsg:
		t.push(msg.toast)
		l.Debug().Int("toast-count", len(t.activeToasts)).Msg("Added toast")
		return t, msg.toast.Init()
	case toastTTLMsg:
		if t.activeToasts != nil {
			_, cmd := msg.toast.Update(msg)
			return t, cmd
		}
		return t, nil
	}
	m, cmd := t.model.Update(msg)
	t.model = m
	return t, cmd
}

func (t toasted) View() string {
	if t.activeToasts == nil || len(t.activeToasts) == 0 {
		return t.model.View()
	}

	stack := &toastStack{toasts: t.activeToasts}

	m := overlay.New(
		stack,
		t.model,
		overlay.Right,
		overlay.Bottom,
		0,
		0,
	)
	return m.View()
}

func WithToast(model tea.Model) tea.Model {
	arr := make([]tea.Model, 0)
	return toasted{model, arr}
}

func Success(msg string) tea.Cmd {
	return func() tea.Msg {
		return viewToastMsg{
			toast: NewSuccess(msg),
		}
	}
}

func Error(msg string) tea.Cmd {
	return func() tea.Msg {
		return viewToastMsg{
			toast: NewError(msg),
		}
	}
}

type viewToastMsg struct {
	toast *Toast
}

type toastTTLMsg struct {
	toast *Toast
}

type toastDeadMsg struct{
	toast *Toast
}

func (t *Toast) Init() tea.Cmd {
	return t.Tick()
}

func (t *Toast) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l := logger.Get()
	switch msg := msg.(type) {
	case toastTTLMsg:
		l.Debug().Dur("toast-ttl", msg.toast.ttl).Msg("Remaining Toast TTL")
		return t, t.Tick()
	}
	return t, nil
}

func (t *Toast) View() string {
	return t.style.Render(t.msg)
}

func (t *Toast) Tick() tea.Cmd {
	return tea.Tick(time.Second, func(ti time.Time) tea.Msg {
		t.ttl -= time.Second
		if t.ttl <= 0 {
			return toastDeadMsg{t}
		}
		return toastTTLMsg{t}
	})
}


type toastStack struct {
	toasts []tea.Model
}

func (t *toastStack) Init() tea.Cmd {
	return nil
}

func (t *toastStack) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *toastStack) View() string {
	res := make([]string, len(t.toasts))
	for _, toast := range t.toasts {
		res = append(res, toast.View())
	}
	return lipgloss.JoinVertical(lipgloss.Bottom, res...)
}
