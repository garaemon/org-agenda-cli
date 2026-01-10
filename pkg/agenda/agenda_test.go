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

func TestAdjustDate(t *testing.T) {
	// Fri, Jan 9, 2026
	baseDate := time.Date(2026, 1, 9, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		rangeType string
		expected  time.Time
	}{
		{
			rangeType: "week",
			// Should be Sun, Jan 4, 2026
			expected: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			rangeType: "month",
			// Should be Jan 1, 2026
			expected: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			rangeType: "day",
			// Should be Jan 9, 2026 (normalized)
			expected: time.Date(2026, 1, 9, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		got := AdjustDate(baseDate, tt.rangeType)
		if !got.Equal(tt.expected) {
			t.Errorf("AdjustDate(%s, %s) = %s; want %s", baseDate, tt.rangeType, got, tt.expected)
		}
	}

	// Test case where today is Sunday
	sunday := time.Date(2026, 1, 4, 10, 0, 0, 0, time.UTC)
	gotSunday := AdjustDate(sunday, "week")
	expectedSunday := time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)
	if !gotSunday.Equal(expectedSunday) {
		t.Errorf("AdjustDate(%s, week) = %s; want %s", sunday, gotSunday, expectedSunday)
	}
}
