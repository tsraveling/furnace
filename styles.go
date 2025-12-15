package main

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("205")
	ColorSecondary = lipgloss.Color("170")
	ColorMuted     = lipgloss.Color("240")

	// Styles

	ViewStyle = lipgloss.NewStyle().
			MarginTop(1).
			PaddingTop(1).
			PaddingLeft(2).
			PaddingBottom(1).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(ColorPrimary)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginLeft(2)

	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(ColorSecondary)

	PaginationStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
