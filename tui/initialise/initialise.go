package initialise

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Inputs       []textinput.Model
	Focused      int
	Done         bool
	UpstreamURL  string
	FirewallPort int
}

var inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))

func New() Model {
	upstream := textinput.New()
	upstream.Placeholder = "http://localhost:8080"
	upstream.Focus()
	upstream.Width = 40

	port := textinput.New()
	port.Placeholder = "5174"
	port.Width = 10
	port.Validate = func(s string) error {
		_, err := strconv.Atoi(s)
		return err
	}

	return Model{
		Inputs: []textinput.Model{upstream, port},
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.Inputs))

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			m.Done = true
			return m, nil

		case "enter":
			// LAST INPUT → SAVE
			if m.Focused == len(m.Inputs)-1 {

				m.UpstreamURL = m.Inputs[0].Value()
				port, _ := strconv.Atoi(m.Inputs[1].Value())
				m.FirewallPort = port

				m.Done = true
				return m, nil
			}

			// otherwise move focus
			m.Focused++
		case "tab":
			m.Focused = (m.Focused + 1) % len(m.Inputs)
		}

		for i := range m.Inputs {
			m.Inputs[i].Blur()
		}
		m.Inputs[m.Focused].Focus()
	}

	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}
func (m Model) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(fmt.Sprintf(
		"Initialization\n\n%s\n%s\n\n%s\n%s\n\nTAB to switch • ESC to save & go back",
		inputStyle.Width(30).Render("Backend URL:"),
		m.Inputs[0].View(),
		inputStyle.Width(30).Render("Firewall Port:"),
		m.Inputs[1].View(),
	))
}
