package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/agenda"
	"github.com/garaemon/org-agenda-cli/pkg/config"
	"github.com/garaemon/org-agenda-cli/pkg/item"
	"github.com/garaemon/org-agenda-cli/pkg/parser"
)

type Service struct {
	OrgFiles    []string
	DefaultFile string
}

func NewService(orgFiles []string, defaultFile string) *Service {
	return &Service{
		OrgFiles:    config.ResolveOrgFiles(orgFiles),
		DefaultFile: defaultFile,
	}
}

type ListOptions struct {
	Status string
	Tag    string
}

func (s *Service) ListTodos(opts ListOptions) ([]*item.Item, error) {
	var allItems []*item.Item
	for _, file := range s.OrgFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			// Skip files that cannot be read, but log/print error?
			// For now, just continue as per original logic
			continue
		}

		items := parser.ParseString(string(content), file)
		for _, it := range items {
			if opts.Status != "" {
				if it.Status != opts.Status {
					continue
				}
			} else {
				if it.Status == "" {
					continue
				}
			}

			if opts.Tag != "" {
				found := false
				for _, t := range it.Tags {
					if t == opts.Tag {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			allItems = append(allItems, it)
		}
	}
	return allItems, nil
}

type AddOptions struct {
	Priority string
	Tags     []string
	Schedule string
	Deadline string
	File     string
}

func (s *Service) AddTodo(title string, opts AddOptions) error {
	targetFile := opts.File
	if targetFile == "" {
		targetFile = s.DefaultFile
	}
	if targetFile == "" && len(s.OrgFiles) > 0 {
		targetFile = s.OrgFiles[0]
	}
	if targetFile == "" {
		return fmt.Errorf("no target file specified")
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.OpenFile(targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	content := fmt.Sprintf("* TODO %s", title)
	if opts.Priority != "" {
		content = fmt.Sprintf("* TODO [#%s] %s", opts.Priority, title)
	}

	if len(opts.Tags) > 0 {
		content += fmt.Sprintf(" :%s:", strings.Join(opts.Tags, ":"))
	}
	content += "\n"

	if opts.Schedule != "" {
		content += fmt.Sprintf("SCHEDULED: <%s>\n", opts.Schedule)
	}
	if opts.Deadline != "" {
		content += fmt.Sprintf("DEADLINE: <%s>\n", opts.Deadline)
	}

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func (s *Service) MarkDone(fileOrId string) error {
	pos, err := parser.ParseFilePosition(fileOrId)
	if err != nil {
		return fmt.Errorf("failed to parse position: %w", err)
	}

	content, err := os.ReadFile(pos.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	if pos.Line < 1 || pos.Line > len(lines) {
		return fmt.Errorf("line number %d out of range", pos.Line)
	}

	lineIdx := pos.Line - 1
	targetLine := lines[lineIdx]

	if !strings.Contains(targetLine, item.StatusTodo) {
		// Try to see if it already is DONE?
		// Original logic: "Line does not appear to be a TODO item"
		return fmt.Errorf("line does not appear to be a %s item", item.StatusTodo)
	}

	newLine := strings.Replace(targetLine, " "+item.StatusTodo+" ", " "+item.StatusDone+" ", 1)
	if newLine == targetLine {
		return fmt.Errorf("could not find ' %s ' pattern to replace", item.StatusTodo)
	}

	lines[lineIdx] = newLine
	newContent := strings.Join(lines, "\n")

	if err := os.WriteFile(pos.FilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *Service) GetAgenda(startedAt time.Time, rangeType string) ([]*item.Item, error) {
	start := agenda.AdjustDate(startedAt, rangeType)
	end := start
	switch rangeType {
	case "week":
		end = start.AddDate(0, 0, 6)
	case "month":
		end = start.AddDate(0, 1, 0)
	default: // day
		// end is same as start (cover the whole day)
		// But FilterItemsByRange uses inclusive comparison.
	}

	// FilterItemsByRange logic seems to handle "day" correctly if start == end?
	// Let's check agenda.FilterItemsByRange implementation again.
	// It truncates to day. So start=2023-01-01, end=2023-01-01 covers that day.

	var allItems []*item.Item
	for _, file := range s.OrgFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		items := parser.ParseString(string(content), file)
		filtered := agenda.FilterItemsByRange(items, start, end)
		allItems = append(allItems, filtered...)
	}

	// Sort by date? The original implementation just appends.
	// Users might appreciate sorting.
	// For now let's keep it simple and consistent with original behavior.

	return allItems, nil
}
