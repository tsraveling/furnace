package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

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

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// SECTION: Core model and view
type pickerModel struct {
	list          list.Model
	input         textinput.Model
	allItems      []list.Item
	hasExactMatch bool
	choice        string
}

func foodPicker() (pickerModel, tea.Cmd) {
	const defaultWidth = 20

	items := cfg.foodDB.All()
	allItems := make([]list.Item, len(items))
	for i, item := range items {
		allItems[i] = item
	}

	lh := min(len(items)+4, listHeight)

	l := list.New(allItems, itemDelegate{}, defaultWidth, lh)
	l.Title = "Item to log?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false)

	ti := textinput.New()
	ti.Placeholder = "Start typing to filter ..."
	ti.Focus()
	ti.CharLimit = 128

	m := pickerModel{list: l, allItems: allItems, input: ti}
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
		m.list.SetWidth(msg.Width)
		m.input.Width = msg.Width

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc", "ctrl+c":
			return m, tea.Quit

		case "ctrl+n":
			if !m.canCreate() {
				break
			}
			// -> Make Item flow
			return makeCreateItemModel(m.input.Value())

		case "enter":
			i, ok := m.list.SelectedItem().(FoodItem)
			if ok {
				m.choice = i.Name
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
	help := helpStyle.Render(help_text)
	return fmt.Sprintf("Food:\n\n%s\n\n%s\n\n%s", m.list.View(), m.input.View(), help)
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginLeft(2)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))

	paginationStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
