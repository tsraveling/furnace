package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type createRecipeField int

const (
	crFieldName         createRecipeField = iota
	crFieldUnits
	crFieldIngredients
	crFieldOtherCalories
)

type createRecipeModel struct {
	quitting         bool
	backState        pickerModel
	focused          createRecipeField
	nameInput        textinput.Model
	unitsInput       textinput.Model
	otherCalInput    textinput.Model
	ingredients      []RecipePart
	ingredientCursor int
	err              error
}

func makeCreateRecipeModel(back pickerModel) (createRecipeModel, tea.Cmd) {
	n := textinput.New()
	n.Placeholder = "e.g. Chicken Stir Fry"
	n.Focus()
	n.Width = min(30, cfg.fullWidth())
	n.Prompt = "> "

	u := textinput.New()
	u.Placeholder = "e.g. bowls, servings"
	u.Width = min(30, cfg.fullWidth())
	u.Prompt = "> "

	oc := textinput.New()
	oc.Placeholder = "e.g. 150 (0 for none)"
	oc.Width = min(30, cfg.fullWidth())
	oc.Prompt = "> "
	oc.Validate = intValidator

	m := createRecipeModel{
		focused:       crFieldName,
		backState:     back,
		nameInput:     n,
		unitsInput:    u,
		otherCalInput: oc,
	}
	return m, m.Init()
}

func (m createRecipeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *createRecipeModel) addIngredient(part RecipePart) {
	m.ingredients = append(m.ingredients, part)
}

func (m *createRecipeModel) refreshFocus() {
	m.nameInput.Blur()
	m.unitsInput.Blur()
	m.otherCalInput.Blur()
	switch m.focused {
	case crFieldName:
		m.nameInput.Focus()
	case crFieldUnits:
		m.unitsInput.Focus()
	case crFieldOtherCalories:
		m.otherCalInput.Focus()
	}
}

func (m *createRecipeModel) checkValid() bool {
	nm := m.nameInput.Value()
	if len(nm) < 1 {
		m.err = fmt.Errorf("Enter a recipe name")
		return false
	}

	if cfg.recipeDB.Get(nm) != nil {
		m.err = fmt.Errorf("\"%s\" already exists as a recipe!", nm)
		return false
	}
	if cfg.foodDB.Get(nm) != nil {
		m.err = fmt.Errorf("\"%s\" already exists as a food item!", nm)
		return false
	}

	if len(m.unitsInput.Value()) < 1 {
		m.err = fmt.Errorf("Enter the units for this recipe")
		return false
	}

	if len(m.ingredients) < 1 {
		m.err = fmt.Errorf("Add at least one ingredient")
		return false
	}

	calErr := intValidator(m.otherCalInput.Value())
	if calErr != nil {
		m.err = fmt.Errorf("Other calories must be a whole number (use 0 for none)")
		return false
	}

	return true
}

func (m createRecipeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 3)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cfg.updateWW(msg.Width)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			return m.backState, m.backState.Init()

		case "tab":
			m.focused = (m.focused + 1) % 4
			m.refreshFocus()
			return m, nil
		case "shift+tab":
			m.focused = (m.focused + 3) % 4
			m.refreshFocus()
			return m, nil

		case "enter":
			switch m.focused {
			case crFieldName, crFieldUnits:
				m.focused++
				m.refreshFocus()
				return m, nil
			case crFieldIngredients:
				m.focused = crFieldOtherCalories
				m.refreshFocus()
				return m, nil
			case crFieldOtherCalories:
				if m.checkValid() {
					otherCal, _ := strconv.Atoi(m.otherCalInput.Value())
					recipe := Recipe{
						Name:          m.nameInput.Value(),
						Units:         m.unitsInput.Value(),
						Parts:         m.ingredients,
						OtherCalories: otherCal,
					}
					cfg.recipeDB.Add(recipe)
					return makeFoodPicker(m.backState.forDate, "", pickerModeLog, nil)
				}
			}

		case "a":
			if m.focused == crFieldIngredients {
				return makeFoodPicker(m.backState.forDate, "", pickerModeIngredient, &m)
			}
		case "x":
			if m.focused == crFieldIngredients && len(m.ingredients) > 0 {
				idx := m.ingredientCursor
				m.ingredients = append(m.ingredients[:idx], m.ingredients[idx+1:]...)
				if m.ingredientCursor >= len(m.ingredients) && m.ingredientCursor > 0 {
					m.ingredientCursor--
				}
				return m, nil
			}
		case "up":
			if m.focused == crFieldIngredients && m.ingredientCursor > 0 {
				m.ingredientCursor--
				return m, nil
			}
		case "down":
			if m.focused == crFieldIngredients && m.ingredientCursor < len(m.ingredients)-1 {
				m.ingredientCursor++
				return m, nil
			}
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.unitsInput, cmds[1] = m.unitsInput.Update(msg)
	m.otherCalInput, cmds[2] = m.otherCalInput.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m createRecipeModel) View() string {
	if m.quitting {
		return quitting()
	}

	title := TitleStyle.Width(cfg.fullWidth()).Render("Create recipe:")

	// Name
	ni := ActiveStyle.Render(m.nameInput.View())

	// Units
	ui := ActiveStyle.Render(m.unitsInput.View())

	// Ingredients section
	var ingredientLines string
	if len(m.ingredients) == 0 {
		ingredientLines = HelpStyle.Render("  (none yet)")
	} else {
		for i, part := range m.ingredients {
			food := cfg.foodDB.Get(part.Food)
			calInfo := ""
			if food != nil {
				calInfo = fmt.Sprintf(" (%d cal)", int(float64(food.Calories)*part.Quantity))
			}
			prefix := "  "
			if m.focused == crFieldIngredients && i == m.ingredientCursor {
				prefix = "● "
			}
			line := fmt.Sprintf("%s%s x%g%s", prefix, part.Food, part.Quantity, calInfo)
			if m.focused == crFieldIngredients && i == m.ingredientCursor {
				ingredientLines += SelectedItemStyle.Render(line) + "\n"
			} else {
				ingredientLines += ItemStyle.Render(line) + "\n"
			}
		}
	}

	var ingredientHelp string
	if m.focused == crFieldIngredients {
		ingredientHelp = HelpStyle.Render("a: add ingredient  x: delete  enter: done")
	}

	// Other calories
	oc := ActiveStyle.Render(m.otherCalInput.View())

	// Error
	errMsg := ""
	if m.err != nil {
		errMsg = ErrorStyle.Render(m.err.Error())
	}

	body := fmt.Sprintf("%s\n\nName:\n%s\n\nUnits:\n%s\n\nIngredients:\n%s%s\n\nOther Calories per unit:\n%s\n\n%s",
		title, ni, ui, ingredientLines, ingredientHelp, oc, errMsg)
	return ViewStyle.Render(body)
}
