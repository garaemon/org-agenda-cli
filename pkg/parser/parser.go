package parser

import (
	"regexp"
	"strings"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/item"
)

var (
	headlineRegex  = regexp.MustCompile(`^\*+\s+(?:(TODO|DONE|WAITING)\s+)?(.*?)(?:\s+:(.*):)?\s*$`)
	timestampRegex = regexp.MustCompile(`<(\d{4}-\d{2}-\d{2})[^>]*>`)
)

// ParseString parses a string containing Org-mode content.
func ParseString(content string, filePath string) []*item.Item {
	lines := strings.Split(content, "\n")
	var items []*item.Item
	var currentItem *item.Item

	for i, line := range lines {
		if strings.HasPrefix(line, "*") {
			currentItem = ParseHeadline(line)
			if currentItem != nil {
				currentItem.FilePath = filePath
				currentItem.LineNumber = i + 1
				items = append(items, currentItem)
			}
			continue
		}

		if currentItem != nil {
			if sched := ParseTimestamp(line, "SCHEDULED"); sched != nil {
				currentItem.Scheduled = sched
			}
			if dead := ParseTimestamp(line, "DEADLINE"); dead != nil {
				currentItem.Deadline = dead
			}
			// For RawContent, we append lines that are not headlines or special metadata
			// In a more complex implementation, we would handle properties etc.
			if !strings.Contains(line, "SCHEDULED:") && !strings.Contains(line, "DEADLINE:") {
				if currentItem.RawContent == "" {
					currentItem.RawContent = line
				} else {
					currentItem.RawContent += "\n" + line
				}
			}
		}
	}

	return items
}

// ParseHeadline parses a single line as an Org headline.
func ParseHeadline(line string) *item.Item {
	matches := headlineRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	level := 0
	for _, char := range line {
		if char == '*' {
			level++
		} else {
			break
		}
	}

	status := matches[1]
	title := matches[2]
	tagsStr := matches[3]

	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ":")
	} else {
		tags = []string{}
	}

	return &item.Item{
		Title:  title,
		Level:  level,
		Status: status,
		Tags:   tags,
	}
}

// ParseTimestamp extracts a timestamp for a given key (e.g., SCHEDULED, DEADLINE).
func ParseTimestamp(line string, key string) *time.Time {
	if !strings.Contains(line, key+":") {
		return nil
	}

	matches := timestampRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return nil
	}

	t, err := time.Parse("2006-01-02", strings.TrimSpace(matches[1]))
	if err != nil {
		return nil
	}

	return &t
}
