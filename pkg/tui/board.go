package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	focusedColumnStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("205")) // Pinkish for focus
)

func ViewBoard(m Model) string {
	var cols []string

	for i := 0; i < 3; i++ {
		// Render each list
		// We might need to adjust size of lists here or in Update?
		// For now, let's assume Update handles sizing (which we need to implement).

		var style lipgloss.Style
		if i == m.focusedInt {
			style = focusedColumnStyle
		} else {
			style = columnStyle
		}

		cols = append(cols, style.Render(m.boardLists[i].View()))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)
}
