package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func quitting() string {
	// Gradient from red to yellow
	furnace := "FURNACE"
	colors := []string{"#FF0000", "#FF3300", "#FF6600", "#FF9900", "#FFCC00", "#FFEE00", "#FFFF00"}
	var styledFurnace string
	for i, c := range furnace {
		styledFurnace += lipgloss.NewStyle().Foreground(lipgloss.Color(colors[i])).Render(string(c))
	}

	ret := "\n" + styledFurnace + "\n\n"

	logs := loadLogs()
	today := time.Now()

	dateStyle := lipgloss.NewStyle().Foreground(ColorBasic).Width(12)
	calStyle := lipgloss.NewStyle().Foreground(ColorPrimary)

	for i := range 5 {
		day := today.AddDate(0, 0, -i)
		y, mo, d := day.Date()
		sum := 0
		for _, l := range logs {
			ly, lm, ld := l.date.Date()
			if ly == y && lm == mo && ld == d {
				sum += l.calories
			}
		}

		date := day.Format("Mon 1.2")
		ds := dateStyle
		cs := calStyle

		if i == 0 {
			ds = ds.Bold(true).Foreground(ColorActive)
			cs = cs.Bold(true)
		}

		ret += ds.Render(date) + cs.Render(fmt.Sprintf("%d", sum)) + "\n"
	}

	return ret
}
