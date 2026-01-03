package tui

import (
	"fmt"
	"time"

	"github.com/Jonathanthedeveloper/havoc.git/internal/logger"
	"github.com/Jonathanthedeveloper/havoc.git/internal/proxy"
	"github.com/Jonathanthedeveloper/havoc.git/internal/state"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "Decrease"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "Increase"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	state    *state.HavocState
	help     help.Model
	keys     keyMap
	quitting bool
	lastKey  string
	latency  progress.Model
	jitter   progress.Model
	dropRate progress.Model
	cursor   int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.latency.Width = msg.Width
		m.jitter.Width = msg.Width
		m.dropRate.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor == 0 {
				m.cursor = 2
			} else {
				m.cursor--
			}
			m.lastKey = "↑"
		case key.Matches(msg, m.keys.Down):
			if m.cursor == 2 {
				m.cursor = 0
			} else {
				m.cursor++
			}
			m.lastKey = "↓"
		case key.Matches(msg, m.keys.Left):
			updateSelectedChaos(m.cursor, m, -1)
			m.lastKey = "←"
		case key.Matches(msg, m.keys.Right):
			updateSelectedChaos(m.cursor, m, 1)
			m.lastKey = "→"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func updateSelectedChaos(cursor int, m model, direction int) {

	choas := m.state.GetChaos()

	switch cursor {
	case 0:
		m.state.SetLatency(choas.Latency + time.Duration(direction)*500*time.Millisecond)
	case 1:
		m.state.SetJitter(choas.Jitter + time.Duration(direction)*100*time.Millisecond)
	case 2:
		m.state.SetDropRate(choas.DropRate + float64(direction)*0.05)
	}
}

func (m model) View() string {

	logo := renderLogo()
	helpView := m.help.View(m.keys)
	configurationView := renderConfiguration(m)
	overviewView := renderOverview(m)

	return logo + "\n\n" + overviewView + "\n\n" + configurationView + "\n\n" + helpView
}

func Start(state *state.HavocState) error {
	// Start the proxy in the background
	go func() {
		if err := proxy.Start(state); err != nil {
			logger.ErrorF("proxy server failed: %s", err)
		}
	}()

	m := model{
		state,
		help.New(),
		keys,
		false,
		"",
		progress.New(progress.WithSolidFill("#FF69B4"), progress.WithoutPercentage()),
		progress.New(progress.WithSolidFill("#FF69B4"), progress.WithoutPercentage()),
		progress.New(progress.WithSolidFill("#FF69B4"), progress.WithoutPercentage()),
		0,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		logger.ErrorF("failed to start %v", err)
		return err
	}

	return nil
}

func renderLogo() string {
	return `
 ░██     ░██    ░███    ░██    ░██   ░██████     ░██████  
░██     ░██   ░██░██   ░██    ░██  ░██   ░██   ░██   ░██ 
░██     ░██  ░██  ░██  ░██    ░██ ░██     ░██ ░██        
░██████████ ░█████████ ░██    ░██ ░██     ░██ ░██        
░██     ░██ ░██    ░██  ░██  ░██  ░██     ░██ ░██        
░██     ░██ ░██    ░██   ░██░██    ░██   ░██   ░██   ░██ 
░██     ░██ ░██    ░██    ░███      ░██████     ░██████
`
}

var cursorBorder = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderLeft(true).BorderForeground(lipgloss.Color("#FF69B4")).PaddingLeft(2)
var normalBorder = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderLeft(true).PaddingLeft(2)

var highlight = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4"))
var normal = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

func applyBorder(view string, cursor bool) string {
	if cursor {
		return cursorBorder.Render(view)
	}
	return normalBorder.Render(view)
}

func applyHighlight(view string, cursor bool) string {
	if cursor {
		return highlight.Render(view)
	}
	return normal.Render(view)
}

func renderConfiguration(m model) string {

	choas := m.state.GetChaos()

	latencyValue := fmt.Sprintf("%v", choas.Latency)
	latencyPercent := float64(choas.Latency) / float64(state.MaxLatency)
	latencyHeader := applyHighlight(lipgloss.JoinHorizontal(lipgloss.Left, "Latency", lipgloss.NewStyle().MarginLeft(1).Render(latencyValue)), m.cursor == 0)
	latencyView := applyBorder(lipgloss.JoinVertical(lipgloss.Left, latencyHeader, m.latency.ViewAs(latencyPercent)), m.cursor == 0)

	jitterValue := fmt.Sprintf("%v", choas.Jitter)
	jitterPercent := float64(choas.Jitter) / float64(state.MaxJitter)
	jitterHeader := applyHighlight(lipgloss.JoinHorizontal(lipgloss.Left, "Jitter", lipgloss.NewStyle().MarginLeft(1).Render(jitterValue)), m.cursor == 1)
	jitterView := applyBorder(lipgloss.JoinVertical(lipgloss.Left, jitterHeader, m.jitter.ViewAs(jitterPercent)), m.cursor == 1)

	dropRateValue := fmt.Sprintf("%.0f%%", choas.DropRate*100)
	dropRatePercent := choas.DropRate / state.MaxDropRate
	dropRateHeader := applyHighlight(lipgloss.JoinHorizontal(lipgloss.Left, "Drop Rate", lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(1).Render(dropRateValue)), m.cursor == 2)
	dropRateView := applyBorder(lipgloss.JoinVertical(lipgloss.Left, dropRateHeader, m.dropRate.ViewAs(dropRatePercent)), m.cursor == 2)

	return lipgloss.JoinVertical(lipgloss.Left, latencyView, "", jitterView, "", dropRateView)
}

func renderOverview(m model) string {
	connection := m.state.GetConnection()

	subtle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	addrStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4")).Bold(true)

	localBlock := lipgloss.JoinHorizontal(lipgloss.Center,
		subtle.Render("[ "),
		labelStyle.Render("FROM"),
		addrStyle.Render(fmt.Sprintf(" :%d ", connection.Port)),
		subtle.Render(" ]"),
	)

	targetBlock := lipgloss.JoinHorizontal(lipgloss.Center,
		subtle.Render("[ "),
		labelStyle.Render("TO"),
		addrStyle.Render(fmt.Sprintf(" %s ", connection.Target)),
		subtle.Render(" ]"),
	)

	arrow := subtle.Render(" ──────▶ ")

	return lipgloss.JoinHorizontal(lipgloss.Center, localBlock, arrow, targetBlock)
}
