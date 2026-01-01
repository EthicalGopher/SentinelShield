package tui

import (
	"encoding/json"
	"os"

	"github.com/EthicalGopher/SentinelShield/tui/all_logs"
	"github.com/EthicalGopher/SentinelShield/tui/initialise"
	"github.com/EthicalGopher/SentinelShield/tui/main_menu"
	tea "github.com/charmbracelet/bubbletea"
)

type View int

const (
	Menu View = iota
	Init
	Logs
)

type Config struct {
	UpstreamURL  string `json:"backend_url"`
	FirewallPort int    `json:"port"`
}
type RootModel struct {
	view View
	menu main_menu.Model
	init initialise.Model
	logs all_logs.Model
	cfg  Config
}

func NewRoot() tea.Model {
	return RootModel{
		view: Menu,
		menu: main_menu.New(),
		init: initialise.New(),
		logs: all_logs.New(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch m.view {

	case Menu:
		{
			var cmd tea.Cmd
			m.menu, cmd = m.menu.Update(msg)

			if m.menu.StartInit {
				m.menu.StartInit = false
				m.view = Init
			}
			if m.menu.StartLogs {
				m.menu.StartLogs = false
				m.view = Logs
			}
			return m, cmd
		}
	case Init:
		{
			var cmd tea.Cmd
			m.init, cmd = m.init.Update(msg)

			if m.init.Done {
				m.cfg.UpstreamURL = m.init.UpstreamURL
				m.cfg.FirewallPort = m.init.FirewallPort

				saveConfig(m.cfg)

				m.init.Done = false
				m.view = Menu
			}
			return m, cmd
		}
	case Logs:
		{
			var cmd tea.Cmd
			m.logs, cmd = m.logs.Update(msg)
			if m.logs.Done {
				m.logs.Done = false
				m.view = Menu
			}
			return m, cmd
		}
	default:
		return m, nil
	}
}
func (m RootModel) View() string {
	if m.view == Init {
		return m.init.View()
	}
	if m.view == Logs {
		return m.logs.View()
	}
	return m.menu.View()

}
func saveConfig(cfg Config) {
	file, err := os.Create("config.json")
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(cfg)
}
