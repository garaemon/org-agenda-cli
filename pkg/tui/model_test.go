package tui

import (
	"testing"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func TestNewModel(t *testing.T) {
	now := time.Now()
	items := []*item.Item{
		{
			Title:     "Test Item",
			Status:    "TODO",
			Scheduled: &now,
			Tags:      []string{"work"},
		},
	}
	m := NewModel(items, "Test Agenda")
	if m.list.Title != "Test Agenda" {
		t.Errorf("expected title 'Test Agenda', got '%s'", m.list.Title)
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
	if desc == "" {
		t.Error("description should not be empty")
	}
}
