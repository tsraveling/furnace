package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type zoomLevel int

const (
	zoomDay zoomLevel = iota
	zoomWeek
	zoomMonth
	zoomYear
)

type summaryViewModel struct {
	table   table.Model
	zoom    zoomLevel
	viewing time.Time
	logs    []log
}

func getColumns() []table.Column {
	w := cfg.fullWidth() - 20
	cal_w := 9
	quan_w := 12
	item_w := max(w-(cal_w+quan_w), 10)
	return []table.Column{
		{Title: "Item", Width: item_w},
		{Title: "Quantity", Width: quan_w},
		{Title: "Calories", Width: cal_w},
	}
}

func makeSummaryViewModel() (summaryViewModel, tea.Cmd) {
	logs := loadLogs()
	rows := rowsForLogs(logs)
	t := table.New(
		table.WithColumns(getColumns()),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(cfg.fullWidth()),
	)
	m := summaryViewModel{logs: logs, viewing: time.Now(), zoom: zoomDay, table: t}
	return m, m.Init()
}

func (m *summaryViewModel) logsForView() []log {
	y, mo, d := m.viewing.Date()
	var filtered []log
	for _, l := range m.logs {
		ly, lm, ld := l.date.Date()
		if ly == y && lm == mo && ld == d {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

func rowsForLogs(l []log) []table.Row {
	ret := make([]table.Row, len(l))
	for i, item := range l {
		ret[i] = table.Row{item.item.Name, fmt.Sprintf("%g %s", item.quantity, item.item.Units), fmt.Sprintf("%d", item.calories)}
	}
	return ret
}

func (m summaryViewModel) Init() tea.Cmd {
	return nil
}

func (m summaryViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cfg.updateWW(msg.Width)
		m.table.SetWidth(cfg.fullWidth())
		m.table.SetColumns(getColumns())
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	// Text input gets the end of it
	// var cmd tea.Cmd
	// m.input, cmd = m.input.Update(msg)
	// return m, cmd

	return m, nil
}

func (m summaryViewModel) View() string {
	// l := m.logsForView()
	table := m.table.View()
	return ViewStyle.Render(table)
}
