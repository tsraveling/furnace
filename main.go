package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

type model struct {
	ww int
	wh int

	// Config options
	homeFolder string
}

func (m *model) readConfig() tea.Model {
	// Get home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// Path to config
	configPath := filepath.Join(homeDir, ".config", "furnace", "config.ini")

	// Load the INI file
	cfg, err := ini.Load(configPath)
	if err != nil {
		panic(err)
	}

	// Read values
	section := cfg.Section("general")
	m.homeFolder = section.Key("homeFolder").String()

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ww = msg.Width
		m.wh = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "hello world! Your home folder is: " + m.homeFolder
}

func main() {
	m := model{}
	m.readConfig()
	p := tea.NewProgram(m)
	p.Run()
}
