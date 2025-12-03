package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	ww int
	wh int
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
	return "hello world"
}

func main() {
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	p.Run()
}
