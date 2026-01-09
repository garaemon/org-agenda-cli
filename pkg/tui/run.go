package tui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func Run(items []*item.Item, start time.Time, viewRange string, title string) error {
	m := NewModel(items, start, viewRange, title)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return nil
}
