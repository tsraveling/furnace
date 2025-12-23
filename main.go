package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// STUB: args can be processed in here

	cfg = readConfig()

	// FIXME: Delete this debug value
	dbgWhich := 0

	var m tea.Model
	switch dbgWhich {

	case 0: // STUB: from arg ``
		m, _ = makeSummaryViewModel(time.Now())
	case 1: // STUB: from arg `log`
		m, _ = makeFoodPicker(time.Now())
	}
	p := tea.NewProgram(m)
	p.Run()
}
