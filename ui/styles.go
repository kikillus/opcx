package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	// Styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Padding(0, 1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Align(lipgloss.Center)

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(special)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))
)
