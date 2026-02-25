package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func TestBoardView(t *testing.T) {
	now := time.Now()
	items := []*item.Item{
		{Title: "Task 1", Status: "TODO"},
		{Title: "Task 2", Status: "WAITING"},
		{Title: "Task 3", Status: "DONE"},
	}
	// Initial viewRange "" means ALL items
	m := NewModel(items, now, "", "Test", "", false)

	// Switch to Board View (List -> SideBySide -> Board)
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updatedModel.(Model) // SideBySide

	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updatedModel.(Model) // Board

	if m.viewMode != ViewModeBoard {
		t.Errorf("Expected ViewModeBoard")
	}

	// Start with Focus 0 (TODO)
	if m.focusedInt != 0 {
		t.Errorf("Expected focus 0")
	}

	// Move right 'l' -> Focus 1 (WAITING)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updatedModel.(Model)

	if m.focusedInt != 1 {
		t.Errorf("Expected focus 1 after moving right")
	}

	// Check item in WAITING list
	if len(m.boardLists[1].Items()) != 1 {
		t.Fatalf("Expected 1 WAITING item, got %d", len(m.boardLists[1].Items()))
	}
	selectedItem := m.boardLists[1].Items()[0].(ListItem).Item
	t.Logf("Selected item before toggle: %s [%s]", selectedItem.Title, selectedItem.Status)

	// Toggle status 't'
	// WAITING -> TODO
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = updatedModel.(Model)

	// Check lists
	// TODO list should have 2 items (Task 1 + Task 2)
	// WAITING list should have 0 items
	// DONE list should have 1 item (Task 3)

	t.Logf("Lists sizes after toggle: TODO=%d WAITING=%d DONE=%d",
		len(m.boardLists[0].Items()),
		len(m.boardLists[1].Items()),
		len(m.boardLists[2].Items()))

	if len(m.boardLists[1].Items()) != 0 {
		t.Errorf("Expected 0 WAITING items, got %d", len(m.boardLists[1].Items()))
	}
	if len(m.boardLists[0].Items()) != 2 {
		t.Errorf("Expected 2 TODO items, got %d", len(m.boardLists[0].Items()))
	}

	// Verify Task 2 is in TODO list with status TODO
	todoItems := m.boardLists[0].Items()
	found := false
	for _, it := range todoItems {
		li := it.(ListItem)
		if li.Item.Title == "Task 2" {
			found = true
			if li.Item.Status != "TODO" {
				t.Errorf("Expected Task 2 status TODO, got %s", li.Item.Status)
			}
			break
		}
	}
	if !found {
		t.Errorf("Task 2 not found in TODO list")
	}
}
