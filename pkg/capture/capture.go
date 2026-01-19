package capture

import (
	"fmt"
	"os"
	"strings"

	"github.com/garaemon/org-agenda-cli/pkg/parser"
)

// Insert appends the entry to the file, respecting the configuration (Heading/OLP).
func Insert(filePath string, heading string, olp []string, entry string, prepend bool) error {
	// Read file
	contentBytes, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	content := string(contentBytes)

	lines := strings.Split(content, "\n")
	// Remove last empty line if it exists (result of trailing newline)
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	// Handle case where file is empty
	if len(lines) == 0 {
		lines = []string{}
	}

	// Detect entry level
	entryLevel := 0
	for _, line := range strings.Split(entry, "\n") {
		item := parser.ParseHeadline(line)
		if item != nil {
			entryLevel = item.Level
			break
		}
	}

	var insertionLine int
	var targetLevel int
	// headingIndex is the line number of the target headline, or -1 if not applicable (root)
	headingIndex := -1

	if len(olp) > 0 {
		// Handle OLP
		headingIndex, insertionLine, targetLevel = findOLPInsertionPoint(lines, olp)
	} else if heading != "" {
		// Handle single heading
		headingIndex, insertionLine, targetLevel = findHeadingInsertionPoint(lines, heading)
	} else {
		// Default: Append to end or prepend
		if prepend {
			insertionLine = 0
		} else {
			insertionLine = len(lines)
		}
		targetLevel = 0
	}

	if insertionLine == -1 {
		return fmt.Errorf("target not found")
	}

	// Refine insertion point for text-only entries (append to immediate body)
	if headingIndex != -1 && entryLevel == 0 && !prepend {
		insertionLine = findEndOfImmediateBody(lines, headingIndex)
	}

	// Calculate actual insertion point based on prepend
	finalInsertionLine := insertionLine
	if (heading != "" || len(olp) > 0) && prepend {
		// If prepending to a heading, insert right after the heading
		if headingIndex != -1 {
			finalInsertionLine = headingIndex + 1
		}
	}

	// Adjust entry level if we are inserting under a heading
	adjustedEntry := entry
	if targetLevel > 0 {
		adjustedEntry = adjustEntryLevel(entry, targetLevel+1)
	}
	adjustedEntry = strings.TrimSuffix(adjustedEntry, "\n")

	// Insert
	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:finalInsertionLine]...)

	newLines = append(newLines, adjustedEntry)
	newLines = append(newLines, lines[finalInsertionLine:]...)

	output := strings.Join(newLines, "\n")
	// Ensure final newline
	if !strings.HasSuffix(output, "\n") {
		output += "\n"
	}

	return os.WriteFile(filePath, []byte(output), 0644)
}

func findHeadingInsertionPoint(lines []string, heading string) (int, int, int) {
	// Scan for * heading
	for i, line := range lines {
		item := parser.ParseHeadline(line)
		if item != nil && item.Title == heading {
			// Found.
			return i, findEndOfSubtree(lines, i, item.Level), item.Level
		}
	}
	return -1, -1, 0
}

func findOLPInsertionPoint(lines []string, olp []string) (int, int, int) {
	scopeStart := -1 // Virtual root
	scopeLevel := 0

	for _, targetTitle := range olp {
		limit := len(lines)
		if scopeStart != -1 {
			limit = findEndOfSubtree(lines, scopeStart, scopeLevel)
		}

		foundIndex := -1
		foundLevel := -1

		startScan := 0
		if scopeStart != -1 {
			startScan = scopeStart + 1
		}

		for i := startScan; i < limit; i++ {
			item := parser.ParseHeadline(lines[i])
			if item != nil {
				if item.Title == targetTitle {
					foundIndex = i
					foundLevel = item.Level
					break
				}
			}
		}

		if foundIndex != -1 {
			scopeStart = foundIndex
			scopeLevel = foundLevel
		} else {
			return -1, -1, 0
		}
	}

	return scopeStart, findEndOfSubtree(lines, scopeStart, scopeLevel), scopeLevel
}

func findEndOfSubtree(lines []string, startIndex int, level int) int {
	for i := startIndex + 1; i < len(lines); i++ {
		item := parser.ParseHeadline(lines[i])
		if item != nil && item.Level <= level {
			return i
		}
	}
	return len(lines)
}

func findEndOfImmediateBody(lines []string, startIndex int) int {
	for i := startIndex + 1; i < len(lines); i++ {
		if parser.ParseHeadline(lines[i]) != nil {
			return i
		}
	}
	return len(lines)
}

func adjustEntryLevel(entry string, targetLevel int) string {
	lines := strings.Split(entry, "\n")
	if len(lines) == 0 {
		return entry
	}

	// Detect initial level of the entry
	initialLevel := 0
	for _, line := range lines {
		item := parser.ParseHeadline(line)
		if item != nil {
			initialLevel = item.Level
			break
		}
	}

	// If no headlines in entry, return as is (it's just body text)
	if initialLevel == 0 {
		return entry
	}

	shift := targetLevel - initialLevel
	if shift == 0 {
		return entry
	}

	var newLines []string
	for _, line := range lines {
		item := parser.ParseHeadline(line)
		if item != nil {
			// Adjust stars
			newLevel := item.Level + shift
			if newLevel < 1 {
				newLevel = 1 // Minimum level 1
			}
			newStars := strings.Repeat("*", newLevel)
			// Replace old stars with new stars
			// We know the line starts with stars.
			// Reconstruct line: newStars + space + rest
			// Need to be careful to preserve other parts.
			// Simple approach: find first space
			firstSpace := strings.Index(line, " ")
			if firstSpace != -1 {
				newLines = append(newLines, newStars+line[firstSpace:])
			} else {
				// Fallback, shouldn't happen for valid headlines
				newLines = append(newLines, newStars+" "+line[item.Level:])
			}
		} else {
			newLines = append(newLines, line)
		}
	}
	return strings.Join(newLines, "\n")
}
