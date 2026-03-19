package main

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("205")
	ColorSecondary = lipgloss.Color("170")
	ColorError     = lipgloss.Color("124")
	ColorMuted     = lipgloss.Color("240")
	ColorDimBright = lipgloss.Color("248")
	ColorBasic     = lipgloss.Color("250")
	ColorActive    = lipgloss.Color("76")
	ColorRecipe    = lipgloss.Color("214")

	// Gradient colors for progress bar — hex because bubbles/progress
	// passes these to go-colorful which requires hex strings.
	GradientGreenDark   = "#00af00" // ANSI 34
	GradientGreenLight  = "#5fd700" // ANSI 76 (ColorActive)
	GradientOrangeDark  = "#af8700" // ANSI 136
	GradientOrangeLight = "#ff8700" // ANSI 208
	GradientRedDark     = "#af0000" // ANSI 124
	GradientRedBright   = "#ff0000" // ANSI 196
	GradientGrayDark    = "#1c1917"
	GradientGrayLight   = "#57534e"
	BarEmptyColor       = "#5f5f00" // ANSI 58

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
			Foreground(ColorPrimary)

	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	ErrorStyle = lipgloss.NewStyle().Foreground(ColorError)

	ActiveStyle = lipgloss.NewStyle().Foreground(ColorActive)

	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(ColorSecondary)

	UnselectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(ColorDimBright)

	RecipeItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(ColorRecipe)

	SelectedRecipeItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(ColorRecipe)

	PaginationStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
