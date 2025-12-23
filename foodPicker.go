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

// SECTION: List delegate
func (i FoodItem) FilterValue() string { return i.Name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(FoodItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Name)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			it := fmt.Sprintf("● %s", strings.Join(s, " "))
			return SelectedItemStyle.Render(it)
		}
	}

	fmt.Fprint(w, fn(str))
}

// SECTION: Core model and view
type pickerModel struct {
	list          list.Model
	input         textinput.Model
	forDate       time.Time
	allItems      []list.Item
	hasExactMatch bool
	choice        string
	ww            int
}

func (m *pickerModel) updateTitle() {
	m.list.Title = "Log an item for " + m.forDate.Format("Mon Jan 2 '06")
}

func makeFoodPicker(t time.Time, ii string) (pickerModel, tea.Cmd) {
	const defaultWidth = 20

	items := cfg.foodDB.All()
	allItems := make([]list.Item, len(items))
	for i, item := range items {
		allItems[i] = item
	}

	lh := min(len(items)+3, listHeight)

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

	m := pickerModel{list: l, allItems: allItems, input: ti, forDate: t}
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
		case "esc", "ctrl+c":
			return m, tea.Quit

		case "ctrl+n":
			if !m.canCreate() {
				break
			}
			// -> Make Item flow
			return makeCreateItemModel(m.input.Value(), m)

		case "enter":
			i, ok := m.list.SelectedItem().(FoodItem)
			if ok {
				m.choice = i.Name
				return makeLogFoodModel(i, m.forDate)
			}
			return m, tea.Quit
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
	title := TitleStyle.Render(m.list.Title)
	help_text := "↑/↓: move • enter: select • esc: quit"
	if len(m.input.Value()) > 0 {
		if m.hasExactMatch {
			help_text += "\n(this item exists)"
		} else {
			help_text += "\nctrl+n: create \"" + m.input.Value() + "\""
		}
	} else {
		help_text += "\nctrl+n: create new item"
	}
	help := HelpStyle.Render(help_text)
	return ViewStyle.Render(fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", title, m.list.View(), m.input.View(), help))
}
