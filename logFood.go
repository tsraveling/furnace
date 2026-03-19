package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type logFoodMode int

const (
	logFoodModeNormal     logFoodMode = iota // log and go to summary
	logFoodModeIngredient                    // capture quantity, return to recipe builder
)

type logFoodModel struct {
	quitting        bool
	forDate         time.Time
	input           textinput.Model
	currentCalories int
	loggingItem     pickerItem
	numericValue    float64
	err             error
	mode            logFoodMode
	backPicker      pickerModel
}

func makeLogFoodModel(pi pickerItem, d time.Time) (logFoodModel, tea.Cmd) {

	ti := textinput.New()
	ti.Placeholder = "# of " + pi.units
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = cfg.fullWidth()
	logs := loadLogs()
	curr := countCaloriesForDate(logs, d)

	m := logFoodModel{input: ti, loggingItem: pi, forDate: d, currentCalories: curr, mode: logFoodModeNormal}
	return m, m.Init()
}

func makeLogFoodModelForIngredient(pi pickerItem, back pickerModel) (logFoodModel, tea.Cmd) {

	ti := textinput.New()
	ti.Placeholder = "# of " + pi.units
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = cfg.fullWidth()

	m := logFoodModel{
		input:           ti,
		loggingItem:     pi,
		mode:            logFoodModeIngredient,
		backPicker:      back,
		currentCalories: -1,
	}
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
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.mode == logFoodModeIngredient && m.backPicker.backRecipe != nil {
				return *m.backPicker.backRecipe, m.backPicker.backRecipe.Init()
			}
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.err == nil && m.numericValue > 0 {
				switch m.mode {
				case logFoodModeNormal:
					if m.loggingItem.isRecipe {
						m.err = writeRecipeLog(m.loggingItem.name, m.numericValue, m.forDate)
					} else {
						m.err = writeLog(m.loggingItem.name, m.numericValue, m.forDate)
					}
					if m.err == nil {
						return makeSummaryViewModel(m.forDate)
					}
				case logFoodModeIngredient:
					part := RecipePart{Food: m.loggingItem.name, Quantity: m.numericValue}
					recipe := m.backPicker.backRecipe
					recipe.addIngredient(part)
					return *recipe, recipe.Init()
				}
			} else if m.numericValue <= 0 {
				m.err = errors.New("Please enter a quantity!")
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// Get the values
	if val, err := strconv.ParseFloat(m.input.Value(), 64); err == nil {
		m.err = nil
		m.numericValue = val
	} else {
		m.err = errors.New("Please enter a valid number!")
	}

	return m, cmd
}

func (m logFoodModel) View() string {
	if m.quitting {
		return quitting()
	}
	title := TitleStyle.Render("Logging " + m.loggingItem.name + ":")
	var helper string
	if len(m.input.Value()) == 0 {
		helper = HelpStyle.Render("Enter a value to see the caloric value.\n")
	} else if m.err != nil {
		helper = ErrorStyle.Render(m.err.Error())
	} else {
		itemAmount := int(float64(m.loggingItem.caloriesPerUnit) * m.numericValue)
		calc := fmt.Sprintf("in %s: %d calories", m.loggingItem.units, itemAmount)
		tot := ""
		if m.currentCalories >= 0 {
			previousAmount := UnselectedItemStyle.Render(strconv.Itoa(m.currentCalories))
			na := m.currentCalories + itemAmount
			var newAmount string
			if na > cfg.dailyTarget {
				newAmount = ErrorStyle.Render(strconv.Itoa(na))
			} else {
				newAmount = ActiveStyle.Render(strconv.Itoa(na))
			}
			targetAmount := UnselectedItemStyle.Render(strconv.Itoa(cfg.dailyTarget))
			tot = fmt.Sprintf("Today: %s -> %s / %s", previousAmount, newAmount, targetAmount)
		}
		helper = fmt.Sprintf("%s\n%s", ActiveStyle.Render(calc), tot)
	}
	body := title + "\n\n" + m.input.View() + "\n\n" + helper
	return ViewStyle.Render(body)
}
