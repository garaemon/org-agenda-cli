package agenda

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/item"
)

// SaveItem updates the item in the file.
// It assumes the item's LineNumber and FilePath are correct.
// It reads the file, finds the line, and replaces it with the new representation of the item.
// Note: This is a simple implementation that replaces the headline line.
// It does not handle multi-line drawer properties updates nicely if we were to change them,
// but for Priority, Status, Tags, it should be fine as they are on the headline.
// For timestamps (SCHEDULED/DEADLINE), we might need to look at next lines.
func SaveItem(it *item.Item) error {
	content, err := os.ReadFile(it.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", it.FilePath, err)
	}

	lines := strings.Split(string(content), "\n")
	if it.LineNumber < 1 || it.LineNumber > len(lines) {
		return fmt.Errorf("line number %d out of range for file %s", it.LineNumber, it.FilePath)
	}

	// 0-indexed line index
	idx := it.LineNumber - 1

	// Reconstruct the headline
	// Format: *... Status [#Priority] Title :Tags:

	// We need to preserve the level (asterisks)
	originalLine := lines[idx]
	level := 0
	for _, char := range originalLine {
		if char == '*' {
			level++
		} else {
			break
		}
	}
	if level == 0 {
		// fallback if line doesn't start with * (shouldn't happen for valid items)
		level = 1
	}

	prefix := strings.Repeat("*", level)

	sb := strings.Builder{}
	sb.WriteString(prefix)

	if it.Status != "" {
		sb.WriteString(" " + it.Status)
	}

	if it.Priority != "" {
		sb.WriteString(fmt.Sprintf(" [#%s]", it.Priority))
	}

	sb.WriteString(" " + it.Title)

	if len(it.Tags) > 0 {
		sb.WriteString(fmt.Sprintf(" :%s:", strings.Join(it.Tags, ":")))
	}

	lines[idx] = sb.String()

	// Update timestamps if changed
	// This is trickier because timestamps are usually on the next line.
	// We check if the next line contains SCHEDULED or DEADLINE and replace/add it.
	// For now, let's implement basic persistence for Headline properties first.
	// If we want to support Date editing, we need to handle it.

	// Handle Date Editing:
	// If we have Scheduled/Deadline, we need to ensure they exist.
	// We traverse next lines until we find a new headline or end of file.
	// If we find SCHEDULED/DEADLINE line, we update it.
	// If not, we insert it after headline.

	// Simplified logic for Phase 1: Just update headline.
	// TODO: Handle Timestamps persistence.

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(it.FilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", it.FilePath, err)
	}

	return nil
}

// UpdateTimestamp updates or inserts SCHEDULED/DEADLINE timestamp.
func UpdateTimestamp(it *item.Item, key string, val time.Time) error {
	content, err := os.ReadFile(it.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", it.FilePath, err)
	}

	lines := strings.Split(string(content), "\n")
	if it.LineNumber < 1 || it.LineNumber > len(lines) {
		return fmt.Errorf("line number %d out of range for file %s", it.LineNumber, it.FilePath)
	}

	// Attempt to find existing timestamp line
	// We search from item's line + 1 until next headline or end of file
	startIdx := it.LineNumber // LineNumber is 1-based, lines is 0-based. item line is lines[LineNumber-1]. So search starts at lines[LineNumber]

	// Check if we are at EOF
	if startIdx >= len(lines) {
		// Just append
		newLine := fmt.Sprintf("%s: <%s>", key, val.Format("2006-01-02 Mon"))
		lines = append(lines, newLine)
	} else {
		found := false

		for i := startIdx; i < len(lines); i++ {
			line := lines[i]
			if strings.HasPrefix(line, "*") {
				// Next headline, stop
				break
			}

			if strings.Contains(line, key+":") {
				// Found it. Replace the timestamp.
				// Pattern: key: <...>
				keyIdx := strings.Index(line, key+":")
				if keyIdx != -1 {
					// Find < after key
					openAngle := strings.Index(line[keyIdx:], "<")
					if openAngle != -1 {
						absoluteOpen := keyIdx + openAngle
						closeAngle := strings.Index(line[absoluteOpen:], ">")
						if closeAngle != -1 {
							absoluteClose := absoluteOpen + closeAngle
							newTs := val.Format("2006-01-02 Mon")
							line = line[:absoluteOpen+1] + newTs + line[absoluteClose:]
							lines[i] = line
							found = true
							break
						}
					}
				}
			}
		}

		if !found {
			// Insert after headline
			newLine := fmt.Sprintf("%s: <%s>", key, val.Format("2006-01-02 Mon"))
			// Insert at startIdx
			lines = append(lines[:startIdx], append([]string{newLine}, lines[startIdx:]...)...)
		}
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(it.FilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", it.FilePath, err)
	}

	return nil
}
