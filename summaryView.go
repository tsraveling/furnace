package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type zoomLevel int

const (
	zoomDay zoomLevel = iota
	zoomWeek
	zoomMonth
	zoomYear
)

type summaryViewModel struct {
	quitting bool
	table    table.Model
	zoom     zoomLevel
	viewing  time.Time
	viewLogs []log
	logs     []log
	total    int
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

func (m *summaryViewModel) reloadLogs() {
	m.logs = loadLogs()
}

func makeSummaryViewModel(d time.Time) (summaryViewModel, tea.Cmd) {
	logs := loadLogs()
	t := table.New(
		table.WithColumns(getColumns()),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(cfg.fullWidth()),
	)
	m := summaryViewModel{logs: logs, viewing: d, zoom: zoomDay, table: t}
	m.refreshRows()
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
		name := item.name
		if item.isRecipe {
			name = "◈ " + name
		}
		ret[i] = table.Row{name, fmt.Sprintf("%g %s", item.quantity, item.units), fmt.Sprintf("%d", item.calories)}
	}
	return ret
}

func (m *summaryViewModel) refreshRows() {
	m.viewLogs = m.logsForView()
	r := rowsForLogs(m.viewLogs)
	m.total = 0
	for _, log := range m.viewLogs {
		m.total += log.calories
	}
	m.table.SetHeight(max(min(len(m.viewLogs)+1, 12), 6))
	m.table.SetRows(r)
}

func (m *summaryViewModel) deleteRow(i int) {
	path := cfg.getPath("logs.md")
	err := deleteLine(path, m.viewLogs[i].line)
	if err != nil {
		panic(err)
	}
	m.reloadLogs()
	m.refreshRows()
	m.table.SetCursor(max(0, m.table.Cursor()-1))
}

func (m summaryViewModel) getProgress() float64 {
	if cfg.dailyTarget <= 0 {
		return 0
	}
	p := float64(m.total) / float64(cfg.dailyTarget)
	if p > 1.0 {
		return 1.0
	}
	return p
}

func (m summaryViewModel) gradientColors() (string, string) {
	p := m.getProgress()
	if p >= 0.9 {
		return GradientRedDark, GradientRedBright
	}
	if p >= 0.7 {
		return GradientOrangeDark, GradientOrangeLight
	}
	return GradientGreenDark, GradientGreenLight
}

func sameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
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
			m.quitting = true
			return m, tea.Quit
		case "a":
			return makeFoodPicker(m.viewing, "", pickerModeLog, nil)
		case "d":
			m.deleteRow(m.table.Cursor())
		case "left", "h":
			m.viewing = m.viewing.Add(-24 * time.Hour)
			m.table.SetCursor(1)
			m.refreshRows()
		case "right", "l":
			m.viewing = m.viewing.Add(24 * time.Hour)
			m.table.SetCursor(1)
			m.refreshRows()
		}
	}

	// Text input gets the end of it
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m summaryViewModel) View() string {
	if m.quitting {
		return quitting()
	}
	now := time.Now()
	var title_text string
	switch m.zoom {
	case zoomDay:
		dfmt := "Monday, Jan 2"
		if sameDay(now, m.viewing) {
			dfmt = "Monday, Jan 2 (Today)"
		} else if m.viewing.Year() != now.Year() {
			dfmt = "Monday, Jan 2, 2006"
		}
		title_text = m.viewing.Format(dfmt)
	}
	title := TitleStyle.AlignHorizontal(lipgloss.Center).Width(cfg.fullWidth()).Render(title_text)
	table := m.table.View()
	w := cfg.fullWidth()

	// Total line: page indicator left, current/target right
	current := ActiveStyle.Render(fmt.Sprintf("%d", m.total))
	target := lipgloss.NewStyle().Foreground(ColorBasic).Render(fmt.Sprintf("%d", cfg.dailyTarget))
	currentTarget := fmt.Sprintf("%s/%s", current, target)
	count := HelpStyle.Render(fmt.Sprintf("%d/%d", m.table.Cursor()+1, len(m.table.Rows())))
	var totalLine string

	if cfg.dailyTarget > 0 {
		// Progress bar
		left, right := m.gradientColors()
		bar := progress.New(progress.WithGradient(left, right), progress.WithoutPercentage())
		spacerWidth := max(w-(lipgloss.Width(currentTarget)+lipgloss.Width(count)), 1)
		bar.Width = spacerWidth - 6
		bar.EmptyColor = BarEmptyColor
		prog := bar.ViewAs(m.getProgress())
		totalLine = lipgloss.JoinHorizontal(lipgloss.Top, count, "  ", prog, "  ", currentTarget)
	} else {
		totalLabel := fmt.Sprintf("Total: %s", ActiveStyle.Render(fmt.Sprintf("%d calories", m.total)))
		spacerWidth := max(w-(lipgloss.Width(totalLabel)+5), 1)
		spacer := strings.Repeat(" ", spacerWidth)
		totalLine = lipgloss.JoinHorizontal(lipgloss.Top, count, spacer, totalLabel)
	}

	//help := "↑↓jk navigate  ←→hl change day  a add  e edit  d delete  E edit food  f fill mode"
	help := "↑↓jk navigate  ←→hl change day  a add  d delete"
	wrappedHelp := wordwrap.String(help, cfg.fullWidth())
	styledHelp := HelpStyle.Render(wrappedHelp)
	body := fmt.Sprintf("%s\n\n%s\n\n%s\n%s", title, table, totalLine, styledHelp)
	return ViewStyle.Render(body)
}
