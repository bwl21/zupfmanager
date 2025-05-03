package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#1E88E5")
	secondaryColor = lipgloss.Color("#FFC107")
	highlightColor = lipgloss.Color("#FF5722")
	textColor      = lipgloss.Color("#FFFFFF")
	subtleColor    = lipgloss.Color("#AAAAAA")
	errorColor     = lipgloss.Color("#F44336")

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			MarginBottom(1)

	HighlightedTextStyle = lipgloss.NewStyle().
				Foreground(highlightColor)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(primaryColor).
				Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	InfoBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)
)
