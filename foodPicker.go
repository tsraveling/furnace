package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const listHeight = 10

// SECTION: Picker item and mode

type pickerMode int

const (
	pickerModeLog        pickerMode = iota // shows foods + recipes, logs to summary
	pickerModeIngredient                   // shows foods only, returns ingredient to recipe builder
)

type pickerItem struct {
	name            string
	units           string
	caloriesPerUnit int
	isRecipe        bool
}

func (p pickerItem) FilterValue() string { return p.name }

// SECTION: List delegate

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(pickerItem)
	if !ok {
		return
	}

	str := i.name
	if i.isRecipe {
		str = "◈ " + str
	}

	if index == m.Index() {
		prefix := fmt.Sprintf("● %s", str)
		if i.isRecipe {
			fmt.Fprint(w, SelectedRecipeItemStyle.Render(prefix))
		} else {
			fmt.Fprint(w, SelectedItemStyle.Render(prefix))
		}
	} else {
		if i.isRecipe {
			fmt.Fprint(w, RecipeItemStyle.Render(str))
		} else {
			fmt.Fprint(w, ItemStyle.Render(str))
		}
	}
}

// SECTION: Core model and view
type pickerModel struct {
	quitting      bool
	list          list.Model
	input         textinput.Model
	forDate       time.Time
	allItems      []list.Item
	hasExactMatch bool
	choice        string
	ww            int
	mode          pickerMode
	backRecipe    *createRecipeModel
}

func (m *pickerModel) updateTitle() {
	if m.mode == pickerModeIngredient {
		m.list.Title = "Select an ingredient"
	} else {
		m.list.Title = "Log an item for " + m.forDate.Format("Mon Jan 2 '06")
	}
}

func makeFoodPicker(t time.Time, ii string, mode pickerMode, backRecipe *createRecipeModel) (pickerModel, tea.Cmd) {
	const defaultWidth = 20

	var allItems []list.Item

	for _, fi := range cfg.foodDB.All() {
		allItems = append(allItems, pickerItem{
			name:            fi.Name,
			units:           fi.Units,
			caloriesPerUnit: fi.Calories,
			isRecipe:        false,
		})
	}

	if mode == pickerModeLog {
		for _, r := range cfg.recipeDB.All() {
			allItems = append(allItems, pickerItem{
				name:            r.Name,
				units:           r.Units,
				caloriesPerUnit: r.CaloriesPerUnit(),
				isRecipe:        true,
			})
		}
	}

	lh := min(len(allItems)+3, listHeight)

	l := list.New(allItems, itemDelegate{}, defaultWidth, lh)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.Styles.PaginationStyle = PaginationStyle
	l.SetShowHelp(false)
	l.SetWidth(cfg.fullWidth())

	ti := textinput.New()
	ti.Placeholder = "Start typing to filter ..."
	ti.Focus()
	ti.CharLimit = 128
	ti.Width = cfg.fullWidth()
	ti.SetValue(ii)

	m := pickerModel{list: l, allItems: allItems, input: ti, forDate: t, mode: mode, backRecipe: backRecipe}
	m.updateTitle()
	return m, m.Init()
}

func (m pickerModel) Init() tea.Cmd {
	return nil
}

func (m *pickerModel) filterList() {
	query := strings.ToLower(m.input.Value())
	if query == "" {
		m.list.SetItems(m.allItems)
		return
	}
	filtered := []list.Item{}
	m.hasExactMatch = false
	for _, item := range m.allItems {
		filterVal := strings.ToLower(item.FilterValue())
		if strings.Contains(filterVal, query) {
			if filterVal == query {
				m.hasExactMatch = true
			}
			filtered = append(filtered, item)
		}
	}
	m.list.SetItems(filtered)
}

func (m *pickerModel) canCreate() bool {
	return !m.hasExactMatch
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cfg.updateWW(msg.Width)
		m.list.SetWidth(cfg.fullWidth())
		m.input.Width = cfg.fullWidth()

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			if m.mode == pickerModeIngredient && m.backRecipe != nil {
				return *m.backRecipe, m.backRecipe.Init()
			}
			m.quitting = true
			return m, tea.Quit

		case "ctrl+n":
			if !m.canCreate() {
				break
			}
			// -> Make Item flow
			return makeCreateItemModel(m.input.Value(), m)

		case "ctrl+r":
			if m.mode != pickerModeLog {
				break
			}
			return makeCreateRecipeModel(m)

		case "enter":
			pi, ok := m.list.SelectedItem().(pickerItem)
			if !ok {
				return m, tea.Quit
			}
			m.choice = pi.name

			if m.mode == pickerModeIngredient {
				return makeLogFoodModelForIngredient(pi, m)
			}
			return makeLogFoodModel(pi, m.forDate)

		case "up", "down":
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	}

	// Text input gets the end of it
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// Filter
	m.filterList()
	return m, cmd
}

func (m pickerModel) View() string {
	if m.quitting {
		return quitting()
	}
	title := TitleStyle.Render(m.list.Title)
	help_text := "↑/↓: move • enter: select • esc: back"
	if len(m.input.Value()) > 0 {
		if m.hasExactMatch {
			help_text += "\n(this item exists)"
		} else {
			help_text += "\nctrl+n: create \"" + m.input.Value() + "\""
		}
	} else {
		help_text += "\nctrl+n: create new item"
	}
	if m.mode == pickerModeLog {
		help_text += "  ctrl+r: create recipe"
	}
	help := HelpStyle.Render(help_text)
	return ViewStyle.Render(fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", title, m.list.View(), m.input.View(), help))
}
