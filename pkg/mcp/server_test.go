package mcp

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/garaemon/org-agenda-cli/pkg/service"
	"github.com/mark3labs/mcp-go/mcp"
)

func createCallToolRequest(name string, args map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      name,
			Arguments: args,
		},
	}
}

func createTempOrgFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test*.org")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()
	t.Cleanup(func() { os.Remove(tmpfile.Name()) })
	return tmpfile.Name()
}

func TestHandleListTodos_All(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\n* DONE Task 2 :work:\n* TODO Task 3 :home:\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("list_todos", map[string]interface{}{})
	result, err := s.handleListTodos(context.Background(), req)
	if err != nil {
		t.Fatalf("handleListTodos returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListTodos returned tool error: %v", result.Content)
	}
}

func TestHandleListTodos_FilterByStatus(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\n* DONE Task 2 :work:\n* TODO Task 3 :home:\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("list_todos", map[string]interface{}{
		"status": "TODO",
	})
	result, err := s.handleListTodos(context.Background(), req)
	if err != nil {
		t.Fatalf("handleListTodos returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListTodos returned tool error: %v", result.Content)
	}
}

func TestHandleListTodos_FilterByTag(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\n* DONE Task 2 :work:\n* TODO Task 3 :home:\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("list_todos", map[string]interface{}{
		"tag": "home",
	})
	result, err := s.handleListTodos(context.Background(), req)
	if err != nil {
		t.Fatalf("handleListTodos returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListTodos returned tool error: %v", result.Content)
	}
}

func TestHandleListTodos_NilArguments(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	// Arguments is nil (no arguments provided)
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "list_todos",
			Arguments: nil,
		},
	}
	result, err := s.handleListTodos(context.Background(), req)
	if err != nil {
		t.Fatalf("handleListTodos returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListTodos returned tool error: %v", result.Content)
	}
}

func TestHandleAddTodo_Success(t *testing.T) {
	filePath := createTempOrgFile(t, "")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("add_todo", map[string]interface{}{
		"title":    "New Task",
		"priority": "A",
		"tags":     "urgent,important",
		"schedule": "2023-10-10",
		"deadline": "2023-10-15",
	})
	result, err := s.handleAddTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("handleAddTodo returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleAddTodo returned tool error: %v", result.Content)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	strContent := string(content)
	if !strings.Contains(strContent, "TODO [#A] New Task :urgent:important:") {
		t.Errorf("File missing expected task header, got: %s", strContent)
	}
	if !strings.Contains(strContent, "SCHEDULED: <2023-10-10>") {
		t.Errorf("File missing schedule, got: %s", strContent)
	}
	if !strings.Contains(strContent, "DEADLINE: <2023-10-15>") {
		t.Errorf("File missing deadline, got: %s", strContent)
	}
}

func TestHandleAddTodo_TitleOnly(t *testing.T) {
	filePath := createTempOrgFile(t, "")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("add_todo", map[string]interface{}{
		"title": "Simple Task",
	})
	result, err := s.handleAddTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("handleAddTodo returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleAddTodo returned tool error: %v", result.Content)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "* TODO Simple Task") {
		t.Errorf("File missing task, got: %s", string(content))
	}
}

func TestHandleAddTodo_MissingTitle(t *testing.T) {
	filePath := createTempOrgFile(t, "")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("add_todo", map[string]interface{}{})
	result, err := s.handleAddTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("handleAddTodo returned unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("Expected tool error when title is missing")
	}
}

func TestHandleAddTodo_InvalidArguments(t *testing.T) {
	filePath := createTempOrgFile(t, "")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	// Arguments is not a map
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "add_todo",
			Arguments: "invalid",
		},
	}
	result, err := s.handleAddTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("handleAddTodo returned unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("Expected tool error for invalid arguments")
	}
}

func TestHandleAddTodo_TargetFile(t *testing.T) {
	targetFile := createTempOrgFile(t, "")
	defaultFile := createTempOrgFile(t, "")
	svc := service.NewService([]string{targetFile, defaultFile}, defaultFile)
	s := NewServer(svc)

	req := createCallToolRequest("add_todo", map[string]interface{}{
		"title": "File-specific Task",
		"file":  targetFile,
	})
	result, err := s.handleAddTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("handleAddTodo returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleAddTodo returned tool error: %v", result.Content)
	}

	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "* TODO File-specific Task") {
		t.Errorf("Target file missing task, got: %s", string(content))
	}
}

func TestHandleMarkDone_Success(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("mark_done", map[string]interface{}{
		"id": filePath + ":1",
	})
	result, err := s.handleMarkDone(context.Background(), req)
	if err != nil {
		t.Fatalf("handleMarkDone returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleMarkDone returned tool error: %v", result.Content)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "* DONE Task 1") {
		t.Errorf("Task not marked as DONE, got: %s", string(content))
	}
}

func TestHandleMarkDone_MissingID(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("mark_done", map[string]interface{}{})
	result, err := s.handleMarkDone(context.Background(), req)
	if err != nil {
		t.Fatalf("handleMarkDone returned unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("Expected tool error when ID is missing")
	}
}

func TestHandleMarkDone_InvalidArguments(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "mark_done",
			Arguments: "invalid",
		},
	}
	result, err := s.handleMarkDone(context.Background(), req)
	if err != nil {
		t.Fatalf("handleMarkDone returned unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("Expected tool error for invalid arguments")
	}
}

func TestHandleGetAgenda_WithDate(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\nSCHEDULED: <2023-10-01>\n* TODO Task 2\nDEADLINE: <2023-10-05>\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("get_agenda", map[string]interface{}{
		"date":  "2023-10-01",
		"range": "day",
	})
	result, err := s.handleGetAgenda(context.Background(), req)
	if err != nil {
		t.Fatalf("handleGetAgenda returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleGetAgenda returned tool error: %v", result.Content)
	}
}

func TestHandleGetAgenda_WeekRange(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\nSCHEDULED: <2023-10-01>\n* TODO Task 2\nDEADLINE: <2023-10-05>\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("get_agenda", map[string]interface{}{
		"date":  "2023-10-01",
		"range": "week",
	})
	result, err := s.handleGetAgenda(context.Background(), req)
	if err != nil {
		t.Fatalf("handleGetAgenda returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleGetAgenda returned tool error: %v", result.Content)
	}
}

func TestHandleGetAgenda_DefaultValues(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\nSCHEDULED: <2023-10-01>\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	// No date and no range specified - should default to today and "day"
	req := createCallToolRequest("get_agenda", map[string]interface{}{})
	result, err := s.handleGetAgenda(context.Background(), req)
	if err != nil {
		t.Fatalf("handleGetAgenda returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleGetAgenda returned tool error: %v", result.Content)
	}
}

func TestHandleGetAgenda_NilArguments(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\nSCHEDULED: <2023-10-01>\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "get_agenda",
			Arguments: nil,
		},
	}
	result, err := s.handleGetAgenda(context.Background(), req)
	if err != nil {
		t.Fatalf("handleGetAgenda returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleGetAgenda returned tool error: %v", result.Content)
	}
}

func TestHandleGetAgenda_InvalidDate(t *testing.T) {
	filePath := createTempOrgFile(t, "* TODO Task 1\n")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	req := createCallToolRequest("get_agenda", map[string]interface{}{
		"date": "not-a-date",
	})
	result, err := s.handleGetAgenda(context.Background(), req)
	if err != nil {
		t.Fatalf("handleGetAgenda returned unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("Expected tool error for invalid date format")
	}
}

func TestHandleAddTodo_TagsParsing(t *testing.T) {
	filePath := createTempOrgFile(t, "")
	svc := service.NewService([]string{filePath}, filePath)
	s := NewServer(svc)

	// Tags with spaces and empty entries
	req := createCallToolRequest("add_todo", map[string]interface{}{
		"title": "Tagged Task",
		"tags":  " tag1 , tag2 , , tag3 ",
	})
	result, err := s.handleAddTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("handleAddTodo returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleAddTodo returned tool error: %v", result.Content)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	strContent := string(content)
	if !strings.Contains(strContent, ":tag1:tag2:tag3:") {
		t.Errorf("Tags not parsed correctly, got: %s", strContent)
	}
}

func TestNewServer(t *testing.T) {
	svc := service.NewService([]string{}, "")
	s := NewServer(svc)
	if s == nil {
		t.Fatal("NewServer returned nil")
	}
	if s.svc != svc {
		t.Error("Server service field not set correctly")
	}
	if s.server == nil {
		t.Error("Server MCPServer field is nil")
	}
}

func TestRegisterTools(t *testing.T) {
	svc := service.NewService([]string{}, "")
	s := NewServer(svc)
	// registerTools should not panic
	s.registerTools()
}
