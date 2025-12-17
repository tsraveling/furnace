package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)

type logFoodModel struct {
	input        textinput.Model
	loggingItem  FoodItem
	numericValue float64
	err          error
}

func makeLogFoodModel(i FoodItem) (logFoodModel, tea.Cmd) {

	ti := textinput.New()
	ti.Placeholder = "# of " + i.Units
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = cfg.fullWidth()

	m := logFoodModel{input: ti, loggingItem: i}
	return m, m.Init()
}

func (m logFoodModel) Init() tea.Cmd {
	return nil
}

func (m logFoodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cfg.updateWW(msg.Width)
		m.input.Width = cfg.fullWidth()

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// Get the values
	if val, err := strconv.ParseFloat(m.input.Value(), 64); err == nil {
		m.err = nil
		m.numericValue = val
	} else {
		m.err = errors.New("please enter a valid number")
	}

	return m, cmd
}

func (m logFoodModel) View() string {
	title := TitleStyle.Render("Logging " + m.loggingItem.Name + ":")
	var helper string
	if len(m.input.Value()) == 0 {
		helper = HelpStyle.Render("Enter a value to see the caloric value.")
	} else if m.err != nil {
		helper = ErrorStyle.Render("Please enter a valid number!")
	} else {
		calc := fmt.Sprintf("in %s: %d calories", m.loggingItem.Units, int(float64(m.loggingItem.Calories)*m.numericValue))
		helper = ActiveStyle.Render(calc)
	}
	body := title + "\n\n" + m.input.View() + "\n\n" + helper
	return ViewStyle.Render(body)
}

// var (
// 	titleStyle = lipgloss.NewStyle().
// 			Bold(true).
// 			Foreground(lipgloss.Color("205")).
// 			MarginLeft(2)
//
// 	itemStyle = lipgloss.NewStyle().
// 			PaddingLeft(4)
// )
