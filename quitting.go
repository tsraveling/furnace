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
	calWidth := 8

	totalCal := 0
	nonZeroDays := 0
	for i := range 5 {
		day := today.AddDate(0, 0, -i)
		sum := countCaloriesForDate(logs, day)
		if sum > 0 {
			totalCal += sum
			nonZeroDays++
		}

		date := day.Format("Mon 1.2")
		ds := dateStyle
		calColor := ColorPrimary
		if sum == 0 {
			calColor = ColorMuted
		}
		cs := lipgloss.NewStyle().Foreground(calColor).Width(calWidth)
		targetColor := ColorMuted
		if sum > cfg.dailyTarget {
			targetColor = ColorError
		}
		ts := lipgloss.NewStyle().Foreground(targetColor)

		if i == 0 {
			ds = ds.Bold(true).Foreground(ColorActive)
			cs = cs.Bold(true)
			ts = ts.Bold(true)
		}

		// Show w/ targets if set, otherwise just show
		if cfg.dailyTarget > 0 {
			ret += ds.Render(date) + cs.Render(fmt.Sprintf("%d", sum)) + ts.Render(fmt.Sprintf("/ %d", cfg.dailyTarget)) + "\n"
		} else {
			ret += ds.Render(date) + cs.Render(fmt.Sprintf("%d", sum)) + "\n"
		}
	}

	avg := 0
	if nonZeroDays > 0 {
		avg = totalCal / nonZeroDays
	}
	color := ColorActive
	if cfg.dailyTarget > 0 && avg > cfg.dailyTarget {
		color = ColorError
	}
	avgCalStyle := lipgloss.NewStyle().Foreground(color).Width(calWidth)
	avgTargetStyle := lipgloss.NewStyle().Foreground(color)
	if cfg.dailyTarget > 0 {
		ret += "\n" + dateStyle.Render("Average") + avgCalStyle.Render(fmt.Sprintf("%d", avg)) + avgTargetStyle.Render(fmt.Sprintf("/ %d", cfg.dailyTarget)) + "\n"
	} else {
		ret += "\n" + dateStyle.Render("Average") + avgCalStyle.Render(fmt.Sprintf("%d", avg)) + "\n"
	}

	return ret
}
