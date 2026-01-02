package agenda

import (
	"testing"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func TestFilterItemsByRange(t *testing.T) {
	d1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	d3 := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)

	items := []*item.Item{
		{Title: "Task 1", Scheduled: &d1},
		{Title: "Task 2", Deadline: &d2},
		{Title: "Task 3", Scheduled: &d3},
	}

	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

	filtered := FilterItemsByRange(items, start, end)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 items, got %d", len(filtered))
	}
}
