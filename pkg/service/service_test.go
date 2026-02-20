package service

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestService_ListTodos(t *testing.T) {
	// Setup
	tmpfile, err := os.CreateTemp("", "test*.org")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := `* TODO Task 1
* DONE Task 2 :work:
* TODO Task 3 :home:
`
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	svc := NewService([]string{tmpfile.Name()}, tmpfile.Name())

	// Test List All TODOs (including DONE)
	items, err := svc.ListTodos(ListOptions{})
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}
	// Check titles
	if items[0].Title != "Task 1" || items[1].Title != "Task 2" || items[2].Title != "Task 3" {
		t.Errorf("Unexpected items: %v", items)
	}

	// Test Filter by Status (TODO)
	items, err = svc.ListTodos(ListOptions{Status: "TODO"})
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 TODO items, got %d", len(items))
	}
	if items[0].Title != "Task 1" || items[1].Title != "Task 3" {
		t.Errorf("Unexpected items: %v", items)
	}

	// Test Filter by Tag
	items, err = svc.ListTodos(ListOptions{Tag: "home"})
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}
	if items[0].Title != "Task 3" {
		t.Errorf("Expected Task 3, got %s", items[0].Title)
	}

	// Test Filter by Status (DONE)
	items, err = svc.ListTodos(ListOptions{Status: "DONE"})
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}
	if items[0].Title != "Task 2" {
		t.Errorf("Expected Task 2, got %s", items[0].Title)
	}
}

func TestService_AddTodo(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.org")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	svc := NewService([]string{}, tmpfile.Name())

	err = svc.AddTodo("New Task", AddOptions{
		Priority: "A",
		Tags:     []string{"urgent"},
		Schedule: "2023-10-10",
	})
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}

	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	strContent := string(content)

	if !strings.Contains(strContent, "* TODO [#A] New Task :urgent:") {
		t.Errorf("File missing task header, got: %s", strContent)
	}
	if !strings.Contains(strContent, "SCHEDULED: <2023-10-10>") {
		t.Errorf("File missing schedule, got: %s", strContent)
	}
}

func TestService_MarkDone(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.org")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := `* TODO Task 1`
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	svc := NewService([]string{tmpfile.Name()}, tmpfile.Name())

	// Mark done using file:line syntax (line 1)
	filePos := tmpfile.Name() + ":1"
	err = svc.MarkDone(filePos)
	if err != nil {
		t.Fatalf("MarkDone failed: %v", err)
	}

	contentBytes, _ := os.ReadFile(tmpfile.Name())
	if !strings.Contains(string(contentBytes), "* DONE Task 1") {
		t.Errorf("Task not marked as DONE, got: %s", string(contentBytes))
	}
}

func TestService_GetAgenda(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.org")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Use a fixed date for testing relative to "today" logic if needed,
	// but here we pass specific dates to GetAgenda.
	// Let's assume today is 2023-10-01.
	content := `* TODO Task 1
SCHEDULED: <2023-10-01>
* TODO Task 2
DEADLINE: <2023-10-05>
* TODO Task 3
SCHEDULED: <2023-11-01>
`
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	svc := NewService([]string{tmpfile.Name()}, tmpfile.Name())

	// Test Day Agenda (2023-10-01)
	date, _ := time.Parse("2006-01-02", "2023-10-01")
	items, err := svc.GetAgenda(date, "day")
	if err != nil {
		t.Fatalf("GetAgenda failed: %v", err)
	}

	// Should include Task 1 (Scheduled today)
	// Task 2 (Deadline is later in the week) should not be included.
	// Task 2 is DEADLINE 2023-10-05.
	// Task 3 is SCHEDULED 2023-11-01.

	if len(items) != 1 {
		t.Errorf("Expected 1 item for day view, got %d", len(items))
	} else if items[0].Title != "Task 1" {
		t.Errorf("Expected Task 1, got %s", items[0].Title)
	}

	// Test Week Agenda (Starts 2023-10-01 which is Sunday)
	// 2023-10-01 is actually a Sunday.
	items, err = svc.GetAgenda(date, "week")
	if err != nil {
		t.Fatalf("GetAgenda failed: %v", err)
	}

	// Should include Task 1 (1st) and Task 2 (5th)
	if len(items) != 2 {
		t.Errorf("Expected 2 items for week view, got %d", len(items))
	}
}
