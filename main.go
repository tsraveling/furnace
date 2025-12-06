package main

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/ini.v1"
)

type Flow int

const (
	AddLog Flow = iota
)

type Node int

const (
	PickFood Node = iota
	EnterQuantity
)

type model struct {
	ww int
	wh int

	// Data
	foodDB *FoodDB

	// View / nav options
	node Node

	// Node states
	pm pickerModel

	// Config options
	homeFolder string
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
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
	m.homeFolder = expandPath(section.Key("homeFolder").String())

	return m
}

func (m model) Init() tea.Cmd {
	return m.pm.Init()
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
	switch m.node {
	case PickFood:
		pm, cmd := m.pm.Update(msg, m.foodDB)
		m.pm = pm
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	switch m.node {
	case PickFood:
		return m.pm.View(m.foodDB)
	}
	return "hello world! Your home folder is: " + m.homeFolder
}

func main() {
	m := model{}
	m.readConfig()

	db, err := LoadFoodDB(filepath.Join(m.homeFolder, "food.md"))
	if err != nil {
		panic(err)
	}
	m.foodDB = db

	p := tea.NewProgram(m)
	p.Run()
}
