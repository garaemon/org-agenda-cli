package parser

import (
	"reflect"
	"testing"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func TestParseFile(t *testing.T) {
	content := `* TODO Task 1 :tag1:
SCHEDULED: <2026-01-01 Thu>
Body of Task 1
* DONE Task 2
DEADLINE: <2026-01-02 Fri>
Body of Task 2
* Non-TODO Headline
Some content
`
	// In a real scenario, we would write this to a temporary file.
	// For now, let's implement ParseString or use a Reader.
	items := ParseString(content, "test.org")

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
		return
	}

	if items[0].Title != "Task 1" || items[0].Status != "TODO" || len(items[0].Tags) != 1 || items[0].Tags[0] != "tag1" {
		t.Errorf("Unexpected item 0: %+v", items[0])
	}
	if items[0].Scheduled == nil || items[0].Scheduled.Format("2006-01-02") != "2026-01-01" {
		t.Errorf("Unexpected scheduled for item 0: %v", items[0].Scheduled)
	}

	if items[1].Title != "Task 2" || items[1].Status != "DONE" {
		t.Errorf("Unexpected item 1: %+v", items[1])
	}
	if items[1].Deadline == nil || items[1].Deadline.Format("2006-01-02") != "2026-01-02" {
		t.Errorf("Unexpected deadline for item 1: %v", items[1].Deadline)
	}

	if items[2].Title != "Non-TODO Headline" || items[2].Status != "" {
		t.Errorf("Unexpected item 2: %+v", items[2])
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected *item.Item
	}{
		{
			name: "Basic TODO item with tags",
			line: "* TODO Task 1 :tag1:tag2:",
			expected: &item.Item{
				Title:  "Task 1",
				Status: "TODO",
				Tags:   []string{"tag1", "tag2"},
			},
		},
		{
			name: "DONE item without tags",
			line: "* DONE Finished task",
			expected: &item.Item{
				Title:  "Finished task",
				Status: "DONE",
				Tags:   []string{},
			},
		},
		{
			name: "Item without status",
			line: "* Just a headline",
			expected: &item.Item{
				Title:  "Just a headline",
				Status: "",
				Tags:   []string{},
			},
		},
		{
			name: "Item with priority A",
			line: "* TODO [#A] Important Task",
			expected: &item.Item{
				Title:    "Important Task",
				Status:   "TODO",
				Priority: "A",
				Tags:     []string{},
			},
		},
		{
			name: "Item with priority B and tags",
			line: "* TODO [#B] Medium Task :tag1:",
			expected: &item.Item{
				Title:    "Medium Task",
				Status:   "TODO",
				Priority: "B",
				Tags:     []string{"tag1"},
			},
		},
		{
			name: "Item with priority C",
			line: "* TODO [#C] Low Task",
			expected: &item.Item{
				Title:    "Low Task",
				Status:   "TODO",
				Priority: "C",
				Tags:     []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseHeadline(tt.line)
			if got.Title != tt.expected.Title {
				t.Errorf("ParseHeadline() Title = %v, want %v", got.Title, tt.expected.Title)
			}
			if got.Status != tt.expected.Status {
				t.Errorf("ParseHeadline() Status = %v, want %v", got.Status, tt.expected.Status)
			}
			if got.Priority != tt.expected.Priority {
				t.Errorf("ParseHeadline() Priority = %v, want %v", got.Priority, tt.expected.Priority)
			}
			if !reflect.DeepEqual(got.Tags, tt.expected.Tags) {
				t.Errorf("ParseHeadline() Tags = %v, want %v", got.Tags, tt.expected.Tags)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		key       string // "SCHEDULED" or "DEADLINE"
		expected  *time.Time
		shouldErr bool
	}{
		{
			name:     "Scheduled timestamp",
			line:     "SCHEDULED: <2026-01-01 Thu>",
			key:      "SCHEDULED",
			expected: ptrTime(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "Deadline timestamp",
			line:     "DEADLINE: <2026-12-31 Wed>",
			key:      "DEADLINE",
			expected: ptrTime(time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTimestamp(tt.line, tt.key)
			if (got == nil) != (tt.expected == nil) {
				t.Errorf("ParseTimestamp() = %v, want %v", got, tt.expected)
				return
			}
			if got != nil && !got.Equal(*tt.expected) {
				t.Errorf("ParseTimestamp() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
