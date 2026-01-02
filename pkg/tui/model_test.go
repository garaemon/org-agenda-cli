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
	m := NewModel(items, "Test Agenda")
	if m.list.Title != "Test Agenda" {
		t.Errorf("expected title 'Test Agenda', got '%s'", m.list.Title)
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
