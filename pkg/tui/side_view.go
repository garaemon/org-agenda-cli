package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func ViewSideBySide(m Model) string {
	// Left pane: List
	// Right pane: Details (Viewport)

	// Layout
	// We assume external Update logic handles content sync and sizing.

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		docStyle.Render(m.list.View()),
		docStyle.Render(m.viewport.View()),
	)
}
