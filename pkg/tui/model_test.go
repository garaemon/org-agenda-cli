package tui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func TestNewModel(t *testing.T) {
	now := time.Now()
	items := []*item.Item{
		{
			Title:      "Test Item",
			Status:     "TODO",
			Scheduled:  &now,
			Tags:       []string{"work"},
			FilePath:   "sample.org",
			LineNumber: 1,
			RawContent: "This is the content.",
		},
	}
	m := NewModel(items, now, "day", "")
	// Title format: Agenda: YYYY-MM-DD - YYYY-MM-DD
	expectedPrefix := "Agenda: "
	if !strings.HasPrefix(m.list.Title, expectedPrefix) {
		t.Errorf("expected title to start with '%s', got '%s'", expectedPrefix, m.list.Title)
	}
	if m.state != listView {
		t.Errorf("expected initial state listView, got %v", m.state)
	}

	// Check list items
	if len(m.list.Items()) != 1 {
		t.Fatalf("expected 1 item, got %d", len(m.list.Items()))
	}

	li, ok := m.list.Items()[0].(ListItem)
	if !ok {
		t.Fatalf("item is not ListItem")
	}

	if li.Title() != "Test Item" {
		t.Errorf("expected item title 'Test Item', got '%s'", li.Title())
	}

	desc := li.Description()
	if !strings.Contains(desc, "sample.org:1") {
		t.Errorf("description should contain file path and line number, got '%s'", desc)
	}

	// Test state transition
	// Simulate window size msg to initialize viewport
	updatedModel, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updatedModel.(Model)

	// Simulate Enter key
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(Model)
	if m.state != detailView {
		t.Errorf("expected state detailView after Enter, got %v", m.state)
	}

	// Check viewport content
	if !strings.Contains(m.viewport.View(), "This is the content.") {
		t.Errorf("viewport should contain item content")
	}

	// Simulate Esc key
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updatedModel.(Model)
	if m.state != listView {
		t.Errorf("expected state listView after Esc, got %v", m.state)
	}
}

func TestPaging(t *testing.T) {
	now := time.Now()
	// Item 1: today
	// Item 2: 7 days later (8th day, should be outside first week: days 0-6)
	item2Date := now.AddDate(0, 0, 7)

	items := []*item.Item{
		{
			Title:     "Item 1",
			Scheduled: &now,
		},
		{
			Title:     "Item 2",
			Scheduled: &item2Date,
		},
	}

	// Initialize with week view
	m := NewModel(items, now, "week", "")

	// Initially, only Item 1 should be visible
	if len(m.list.Items()) != 1 {
		t.Errorf("Expected 1 item initially, got %d", len(m.list.Items()))
	}
	if m.list.Items()[0].(ListItem).Item.Title != "Item 1" {
		t.Errorf("Expected Item 1, got %s", m.list.Items()[0].(ListItem).Item.Title)
	}

	// Press 'n' to go to next week
	// KeyMsg for 'n' (rune)
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = updatedModel.(Model)

	// Now Item 2 should be visible, Item 1 should be gone
	if len(m.list.Items()) != 1 {
		t.Fatalf("Expected 1 item after paging next, got %d", len(m.list.Items()))
	}
	if m.list.Items()[0].(ListItem).Item.Title != "Item 2" {
		t.Errorf("Expected Item 2, got %s", m.list.Items()[0].(ListItem).Item.Title)
	}

	// Press 'p' to go back
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	m = updatedModel.(Model)

	// Back to Item 1
	if len(m.list.Items()) != 1 {
		t.Fatalf("Expected 1 item after paging back, got %d", len(m.list.Items()))
	}
	if m.list.Items()[0].(ListItem).Item.Title != "Item 1" {
		t.Errorf("Expected Item 1, got %s", m.list.Items()[0].(ListItem).Item.Title)
	}
}

func TestPagingMonth(t *testing.T) {
	now := time.Now()
	// Item 1: today
	// Item 2: 32 days later (outside first month)
	item2Date := now.AddDate(0, 0, 32)

	items := []*item.Item{
		{
			Title:     "Item 1",
			Scheduled: &now,
		},
		{
			Title:     "Item 2",
			Scheduled: &item2Date,
		},
	}

	// Initialize with month view
	m := NewModel(items, now, "month", "")

	// Initially, only Item 1 should be visible
	if len(m.list.Items()) != 1 {
		t.Errorf("Expected 1 item initially, got %d", len(m.list.Items()))
	}
	if m.list.Items()[0].(ListItem).Item.Title != "Item 1" {
		t.Errorf("Expected Item 1, got %s", m.list.Items()[0].(ListItem).Item.Title)
	}

	// Press 'n' to go to next month
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = updatedModel.(Model)

	// Now Item 2 should be visible, Item 1 should be gone
	if len(m.list.Items()) != 1 {
		t.Fatalf("Expected 1 item after paging next, got %d", len(m.list.Items()))
	}
	if m.list.Items()[0].(ListItem).Item.Title != "Item 2" {
		t.Errorf("Expected Item 2, got %s", m.list.Items()[0].(ListItem).Item.Title)
	}

	// Press 'p' to go back
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	m = updatedModel.(Model)

	// Back to Item 1
	if len(m.list.Items()) != 1 {
		t.Fatalf("Expected 1 item after paging back, got %d", len(m.list.Items()))
	}
	if m.list.Items()[0].(ListItem).Item.Title != "Item 1" {
		t.Errorf("Expected Item 1, got %s", m.list.Items()[0].(ListItem).Item.Title)
	}
}
