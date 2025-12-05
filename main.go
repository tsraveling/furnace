package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

type model struct {
	ww int
	wh int
}

func (m model) readConfig() {
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
	fmt.Println("Reading test section of config.ini:")
	section := cfg.Section("test")
	value := section.Key("someValue").String()
	fmt.Println(value)

	// Or with a default value
	// port := section.Key("port").MustInt(8080)
	// fmt.Println(port)
}

func (m model) Init() tea.Cmd {
	m.readConfig()
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
	return "hello world"
}

func main() {
	p := tea.NewProgram(model{})
	p.Run()
}
