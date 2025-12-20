package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// STUB: args can be processed in here

	cfg = readConfig()

	m, _ := foodPicker()
	p := tea.NewProgram(m)
	p.Run()
}
