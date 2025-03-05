package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Message types for the UI
type (
	startTrackingMsg struct{}
	eventMsg        string
	statsMsg        map[string]int64
	errorMsg        error
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A7A7A7")).
		Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Padding(0, 1)

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Padding(0, 1)
)

type model struct {
	spinner     spinner.Model
	list        list.Model
	viewport    viewport.Model
	textinput   textinput.Model
	table       table.Model
	events      []string
	username    string
	stats       map[string]int64
	err         error
	loading     bool
	showList    bool
	showViewer  bool
}

func NewModel(username string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Live Events"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	v := viewport.New(80, 24)
	v.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1)

	ti := textinput.New()
	ti.Placeholder = "Enter username..."
	ti.Focus()

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Metric", Width: 20},
			{Title: "Value", Width: 20},
		}),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	return model{
		spinner:   s,
		list:      l,
		viewport:  v,
		textinput: ti,
		table:     t,
		username:  username,
		stats:     make(map[string]int64),
		showViewer: true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		textinput.Blink,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.viewport.Width = msg.Width - h
		m.viewport.Height = msg.Height - v
		m.list.SetWidth(msg.Width - h)
		m.list.SetHeight(msg.Height - v)
	case eventMsg:
		m.events = append(m.events, string(msg))
		m.viewport.SetContent(m.formatEvents())
		m.viewport.GotoBottom()
	case statsMsg:
		m.stats = map[string]int64(msg)
		m.updateStats()
	case errorMsg:
		m.err = error(msg)
	}

	if m.showViewer {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.showList {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.textinput, cmd = m.textinput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.showViewer {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			titleStyle.Render(fmt.Sprintf("Live Stream: @%s", m.username)),
			m.table.View(),
			m.viewport.View(),
		)
	}

	if m.showList {
		return m.list.View()
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		titleStyle.Render("TikTok Live Logger"),
		infoStyle.Render("Enter a TikTok username to start tracking their live stream:"),
		m.textinput.View(),
	)
}

func (m *model) updateStats() {
	rows := []table.Row{
		{"Viewers", fmt.Sprintf("%d", m.stats["viewers"])},
		{"Likes", fmt.Sprintf("%d", m.stats["likes"])},
		{"Shares", fmt.Sprintf("%d", m.stats["shares"])},
		{"Comments", fmt.Sprintf("%d", m.stats["comments"])},
	}
	m.table.SetRows(rows)
}

func (m *model) formatEvents() string {
	var s string
	for _, event := range m.events {
		s += fmt.Sprintf("%s\n", event)
	}
	return s
}

func (m *model) AddEvent(event string) {
	m.events = append(m.events, event)
	m.viewport.SetContent(m.formatEvents())
	m.viewport.GotoBottom()
}

func (m *model) UpdateStats(stats map[string]int64) {
	m.stats = stats
	m.updateStats()
}

func (m *model) SetError(err error) {
	m.err = err
	m.loading = false
}

func (m *model) SetListItems(items []list.Item) {
	m.list.SetItems(items)
}

func (m *model) SetViewMode(showList bool) {
	m.showList = showList
} 