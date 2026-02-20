package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func TestSorting(t *testing.T) {
	now := time.Now()
	items := []*item.Item{
		{Title: "A", Priority: "A"}, // Priority A
		{Title: "C", Priority: "C"}, // Priority C
		{Title: "B", Priority: "B"}, // Priority B
	}
	// Initial: default order (file order: A, C, B)
	m := NewModel(items, now, "", "Test", "", false)

	// file order check: A, C, B
	if m.list.Items()[0].(ListItem).Item.Title != "A" {
		t.Errorf("Expected A first initially")
	}
	if m.list.Items()[1].(ListItem).Item.Title != "C" {
		t.Errorf("Expected C second initially")
	}
	if m.list.Items()[2].(ListItem).Item.Title != "B" {
		t.Errorf("Expected B third initially")
	}

	// Press 'S' -> Sort by Priority (Asc: A, B, C)
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'S'}})
	m = updatedModel.(Model)

	if m.sortBy != "priority" {
		t.Errorf("Expected sortBy 'priority', got '%s'", m.sortBy)
	}
	if m.list.Items()[0].(ListItem).Item.Title != "A" {
		t.Errorf("Expected A first (priority asc)")
	}
	if m.list.Items()[1].(ListItem).Item.Title != "B" {
		t.Errorf("Expected B second (priority asc), got %s", m.list.Items()[1].(ListItem).Item.Title)
	}
	if m.list.Items()[2].(ListItem).Item.Title != "C" {
		t.Errorf("Expected C third (priority asc)")
	}

	// Press 'O' -> Toggle Desc (Desc: C, B, A)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'O'}})
	m = updatedModel.(Model)

	if !m.sortDesc {
		t.Errorf("Expected sortDesc true")
	}
	if m.list.Items()[0].(ListItem).Item.Title != "C" {
		t.Errorf("Expected C first (priority desc), got %s", m.list.Items()[0].(ListItem).Item.Title)
	}
	if m.list.Items()[2].(ListItem).Item.Title != "A" {
		t.Errorf("Expected A third (priority desc)")
	}

	// Press 'S' again -> Date
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'S'}})
	m = updatedModel.(Model)
	if m.sortBy != "date" {
		t.Errorf("Expected sortBy 'date', got '%s'", m.sortBy)
	}
}
