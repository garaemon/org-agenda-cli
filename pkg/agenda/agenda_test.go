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

	// Test boundary conditions
	dStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	dEnd := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)
	dJustAfterEnd := time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC)

	boundaryItems := []*item.Item{
		{Title: "Start", Scheduled: &dStart},
		{Title: "End", Scheduled: &dEnd},
		{Title: "JustAfter", Scheduled: &dJustAfterEnd},
	}

	res := FilterItemsByRange(boundaryItems, start, end)
	if len(res) != 2 {
		t.Errorf("Expected 2 boundary items, got %d", len(res))
	}
}

func TestExtractUniqueTags(t *testing.T) {
	items := []*item.Item{
		{Title: "Task 1", Tags: []string{"tag1", "tag2"}},
		{Title: "Task 2", Tags: []string{"tag2", "tag3"}},
		{Title: "Task 3", Tags: []string{"tag1"}},
	}

	tags := ExtractUniqueTags(items)

	if len(tags) != 3 {
		t.Errorf("Expected 3 unique tags, got %d", len(tags))
	}

	expected := map[string]bool{"tag1": true, "tag2": true, "tag3": true}
	for _, tag := range tags {
		if !expected[tag] {
			t.Errorf("Unexpected tag: %s", tag)
		}
	}
}
