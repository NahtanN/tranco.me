package styles

import "github.com/charmbracelet/lipgloss"

func Red(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(text)
}

func Green(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(text)
}

func Yellow(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render(text)
}
