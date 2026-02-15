package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/service"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	svc    *service.Service
	server *server.MCPServer
}

func NewServer(svc *service.Service) *Server {
	s := server.NewMCPServer("org-agenda-mcp", "0.1.0")
	return &Server{
		svc:    svc,
		server: s,
	}
}

func (s *Server) Start() error {
	s.registerTools()
	return server.ServeStdio(s.server)
}

func (s *Server) registerTools() {
	s.server.AddTool(mcp.NewTool("list_todos",
		mcp.WithDescription("List TODO items with optional filtering"),
		mcp.WithString("status",
			mcp.Description("Filter by status (TODO, DONE, WAITING)"),
		),
		mcp.WithString("tag",
			mcp.Description("Filter by tag"),
		),
	), s.handleListTodos)

	s.server.AddTool(mcp.NewTool("add_todo",
		mcp.WithDescription("Add a new TODO item"),
		mcp.WithString("title",
			mcp.Description("Title of the task"),
			mcp.Required(),
		),
		mcp.WithString("file",
			mcp.Description("Target file path (optional)"),
		),
		mcp.WithString("priority",
			mcp.Description("Priority (A, B, C)"),
		),
		mcp.WithString("tags",
			mcp.Description("Comma-separated tags"),
		),
		mcp.WithString("schedule",
			mcp.Description("Scheduled date (YYYY-MM-DD or <YYYY-MM-DD ...>)"),
		),
		mcp.WithString("deadline",
			mcp.Description("Deadline date (YYYY-MM-DD or <YYYY-MM-DD ...>)"),
		),
	), s.handleAddTodo)

	s.server.AddTool(mcp.NewTool("mark_done",
		mcp.WithDescription("Mark a task as DONE"),
		mcp.WithString("id",
			mcp.Description("Task ID (currently file path:line number, e.g., /path/to/file.org:10)"),
			mcp.Required(),
		),
	), s.handleMarkDone)

	s.server.AddTool(mcp.NewTool("get_agenda",
		mcp.WithDescription("Get agenda items for a specific date range"),
		mcp.WithString("date",
			mcp.Description("Reference date (YYYY-MM-DD)"),
		),
		mcp.WithString("range",
			mcp.Description("Range type (day, week, month)"),
		),
	), s.handleGetAgenda)
}

func (s *Server) handleListTodos(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		// If no arguments provided, it might be nil or empty map?
		// We can proceed with empty args
		args = make(map[string]interface{})
	}

	status, _ := args["status"].(string)
	tag, _ := args["tag"].(string)

	items, err := s.svc.ListTodos(service.ListOptions{
		Status: status,
		Tag:    tag,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list todos: %v", err)), nil
	}

	return mcp.NewToolResultJSON(items)
}

func (s *Server) handleAddTodo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}

	title, ok := args["title"].(string)
	if !ok {
		return mcp.NewToolResultError("Title is required"), nil
	}
	file, _ := args["file"].(string)
	priority, _ := args["priority"].(string)
	tagsStr, _ := args["tags"].(string)
	schedule, _ := args["schedule"].(string)
	deadline, _ := args["deadline"].(string)

	var tags []string
	if tagsStr != "" {
		for _, t := range strings.Split(tagsStr, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	err := s.svc.AddTodo(title, service.AddOptions{
		File:     file,
		Priority: priority,
		Tags:     tags,
		Schedule: schedule,
		Deadline: deadline,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to add todo: %v", err)), nil
	}

	return mcp.NewToolResultText("Todo added successfully"), nil
}

func (s *Server) handleMarkDone(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}

	id, ok := args["id"].(string)
	if !ok {
		return mcp.NewToolResultError("ID is required"), nil
	}

	err := s.svc.MarkDone(id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to mark done: %v", err)), nil
	}

	return mcp.NewToolResultText("Task marked as DONE"), nil
}

func (s *Server) handleGetAgenda(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		args = make(map[string]interface{})
	}

	dateStr, _ := args["date"].(string)
	rangeType, _ := args["range"].(string)

	date := time.Now()
	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid date format: %v", err)), nil
		}
		date = parsed
	}

	if rangeType == "" {
		rangeType = "day"
	}

	items, err := s.svc.GetAgenda(date, rangeType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get agenda: %v", err)), nil
	}

	return mcp.NewToolResultJSON(items)
}
