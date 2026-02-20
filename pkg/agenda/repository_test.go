package agenda

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func TestSaveItem(t *testing.T) {
	// Create a temporary file
	content := `* TODO [#A] Test Item :work:
Body content
** Another item`
	tmpfile, err := os.CreateTemp("", "agenda_test_*.org")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Create item representing the first line
	it := &item.Item{
		Title:      "Updated Item",
		Status:     "DONE",
		Priority:   "B",
		Tags:       []string{"home"},
		FilePath:   tmpfile.Name(),
		LineNumber: 1,
	}

	if err := SaveItem(it); err != nil {
		t.Fatalf("SaveItem failed: %v", err)
	}

	// Read back
	newContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(newContent), "\n")
	expected := `* DONE [#B] Updated Item :home:`
	if lines[0] != expected {
		t.Errorf("Expected line 0 to be '%s', got '%s'", expected, lines[0])
	}
	if lines[1] != "Body content" {
		t.Errorf("Expected body unchanged")
	}
}

func TestUpdateTimestamp(t *testing.T) {
	// Create a temporary file
	content := `* TODO Test Timestamp
SCHEDULED: <2023-01-01 Sun>
Body`
	tmpfile, err := os.CreateTemp("", "agenda_ts_test_*.org")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	it := &item.Item{
		FilePath:   tmpfile.Name(),
		LineNumber: 1,
	}

	newDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	if err := UpdateTimestamp(it, "SCHEDULED", newDate); err != nil {
		t.Fatalf("UpdateTimestamp failed: %v", err)
	}

	// Read back
	newContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Expect SCHEDULED to be updated
	expectedTs := "SCHEDULED: <2025-12-31 Wed>"
	if !strings.Contains(string(newContent), expectedTs) {
		t.Errorf("Expected content to contain '%s', got:\n%s", expectedTs, string(newContent))
	}

	// Test adding DEADLINE (new)
	deadlineDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := UpdateTimestamp(it, "DEADLINE", deadlineDate); err != nil {
		t.Fatalf("UpdateTimestamp for Deadline failed: %v", err)
	}

	newContent, _ = os.ReadFile(tmpfile.Name())
	expectedDl := "DEADLINE: <2026-01-01 Thu>"
	if !strings.Contains(string(newContent), expectedDl) {
		t.Errorf("Expected content to contain '%s', got:\n%s", expectedDl, string(newContent))
	}
}
