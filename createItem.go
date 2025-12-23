package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)

type createItemModel struct {
	backState  pickerModel
	focused    int
	nameInput  textinput.Model
	unitsInput textinput.Model
	calInput   textinput.Model
	err        error
}

// nextInput focuses the next input field
func (m *createItemModel) nextInput() {
	m.focused = (m.focused + 1) % 3
}

// prevInput focuses the previous input field
func (m *createItemModel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = 2
	}
}

func intValidator(s string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func (m *createItemModel) checkValid() bool {
	nm := m.nameInput.Value()
	if len(nm) < 1 {
		m.err = fmt.Errorf("Enter a name")
		return false
	}

	_, ok := cfg.foodDB.byName[nm]
	if ok {
		m.err = fmt.Errorf("\"%s\" already exists in the database!", nm)
		return false
	}

	if len(m.unitsInput.Value()) < 1 {
		m.err = fmt.Errorf("Enter a the units for this item")
		return false
	}

	calErr := intValidator(m.calInput.Value())
	if calErr != nil {
		m.err = fmt.Errorf("Calories must be a whole number")
		return false
	}

	return true
}

func makeCreateItemModel(nn string, back pickerModel) (createItemModel, tea.Cmd) {

	n := textinput.New()
	n.Placeholder = "e.g. Ramen"
	n.Focus()
	n.Width = min(30, cfg.fullWidth())
	n.Prompt = "> "
	n.SetValue(nn)

	u := textinput.New()
	u.Placeholder = "e.g. servings, tbsp, cups"
	u.Width = min(30, cfg.fullWidth())
	u.Prompt = "> "

	c := textinput.New()
	c.Placeholder = "e.g. 150"
	c.Width = min(30, cfg.fullWidth())
	c.Prompt = "> "
	c.Validate = intValidator

	m := createItemModel{focused: 0, backState: back, nameInput: n, unitsInput: u, calInput: c}
	return m, m.Init()
}

func (m createItemModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *createItemModel) refreshFocus() {
	m.nameInput.Blur()
	m.unitsInput.Blur()
	m.calInput.Blur()
	switch m.focused {
	case 0:
		m.nameInput.Focus()
	case 1:
		m.unitsInput.Focus()
	case 2:
		m.calInput.Focus()
	}
}

func (m createItemModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 3)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cfg.updateWW(msg.Width)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m.backState, m.backState.Init()
		case "shift+tab":
			m.prevInput()
			m.refreshFocus()
		case "tab":
			m.nextInput()
			m.refreshFocus()
		case "enter":
			if m.focused < 2 {
				m.nextInput()
				m.refreshFocus()
				return m, nil
			} else {
				if m.checkValid() {
					// If the new item is valid, save it and go to quantity entry
					cals, err := strconv.ParseInt(m.calInput.Value(), 10, 64)
					if err != nil {
						panic(err) // just in case validation misses it somehow
					}
					fi := FoodItem{Name: m.nameInput.Value(), Units: m.unitsInput.Value(), Calories: int(cals)}
					cfg.foodDB.Add(fi)
					return makeLogFoodModel(fi, m.backState.forDate)
				}
			}
		}

	case error:
		m.err = msg
		return m, nil
	}

	// Process text inputs
	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.unitsInput, cmds[1] = m.unitsInput.Update(msg)
	m.calInput, cmds[2] = m.calInput.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m createItemModel) View() string {
	title := TitleStyle.Width(cfg.fullWidth()).Render("Create item:")
	ni := ActiveStyle.Render(m.nameInput.View())
	ui := ActiveStyle.Render(m.unitsInput.View())
	ci := ActiveStyle.Render(m.calInput.View())
	errMsg := ""
	if m.err != nil {
		errMsg = ErrorStyle.Render(m.err.Error())
	}
	body := fmt.Sprintf("%s\n\nName:\n%s\n\nUnits:\n%s\n\nCalories:\n%s\n\n%s", title, ni, ui, ci, errMsg)
	return ViewStyle.Render(body)
}
