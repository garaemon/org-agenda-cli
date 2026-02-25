package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func TestSideBySide(t *testing.T) {
	now := time.Now()
	items := []*item.Item{
		{Title: "Task 1", RawContent: "Details 1"},
	}
	m := NewModel(items, now, "", "Test", "", false)

	// Set initial size
	m.width = 100
	m.height = 40
	m.list.SetSize(100, 40)

	// Initial View should not be side-by-side
	if m.viewMode != ViewModeList {
		t.Error("Expected ViewModeList initially")
	}

	// Press Tab to toggle
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updatedModel.(Model)

	if m.viewMode != ViewModeSideBySide {
		t.Error("Expected ViewModeSideBySide after Tab")
	}

	// View should contain side-by-side elements (lipgloss JoinHorizontal)
	// We can't easily check rendered string structure but we can check if View() runs without panic
	view := m.View()
	if len(view) == 0 {
		t.Error("Expected non-empty view")
	}

	// Check if viewport content is set
	// In Update loop, we set viewport content if sideBySide is true
	// Wait, we need to trigger update again or did logic run?
	// The logic for updating viewport content is in `case listView:` block in Update method.
	// But `KeyTab` is handled in `KeyMsg` switch, and then it returns `m, nil`.
	// Use of `KeyTab` returns `m, nil` immediately inside the switch.
	// So it doesn't fall through to the end where `listView` logic is?
	// Wait, `Update` logic structure:
	/*
	   switch msg := msg.(type) {
	     case tea.KeyMsg:
	        // handle keys
	        return m, nil // RETURNS HERE
	   }

	   switch m.state {
	     case listView:
	        // update viewport content logic is HERE
	   }
	*/
	// So `KeyTab` updates `m.sideBySide` but DOES NOT update viewport content immediately.
	// However, `View()` calls `ViewSideBySide(m)`.
	// `ViewSideBySide` uses `m.viewport.View()`.
	// `m.viewport` content is not set yet?
	// `m.viewport` was initialized with empty content?
	// But `m.list.SelectedItem()` exists.

	// Issue: Viewport content needs to be initialized.
	// When we switch to side-by-side, we should probably update viewport content immediately.
	// I should update `pkg/tui/model.go` to update viewport content when toggling.

	// Let's verifying this test failure first?
	// Or just fix the code.

	// If I fix the code, I should do it now.
}
