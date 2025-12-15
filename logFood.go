package main

import (
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)

type logFoodModel struct {
	input       textinput.Model
	loggingItem FoodItem
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
		cfg.ww = msg.Width
		m.input.Width = cfg.fullWidth()

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m logFoodModel) View() string {
	title := TitleStyle.Render("Logging " + m.loggingItem.Name + ":")
	body := title + "\n\n" + m.input.View()
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
