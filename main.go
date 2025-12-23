package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	// Load the config file
	cfg = readConfig()

	var m tea.Model
	if len(os.Args) < 2 {
		m, _ = makeSummaryViewModel(time.Now())
	} else {
		switch os.Args[1] {
		case "log":
			// `furn log`
			input := strings.Join(os.Args[2:], " ")
			m, _ = makeFoodPicker(time.Now(), input)
		default:
			fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
			os.Exit(1)
		}
	}

	p := tea.NewProgram(m)
	p.Run()
}
