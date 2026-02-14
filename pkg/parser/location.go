package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// FilePosition represents a specific line in a file.
type FilePosition struct {
	FilePath string
	Line     int
}

// ParseFilePosition parses a string in the format "filepath:line" and returns a FilePosition.
func ParseFilePosition(arg string) (*FilePosition, error) {
	parts := strings.Split(arg, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format: %s. Use 'file:line'", arg)
	}

	filePath := parts[0]
	if filePath == "" {
		return nil, fmt.Errorf("empty file path")
	}

	lineIdx, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid line number: %v", err)
	}

	if lineIdx < 1 {
		return nil, fmt.Errorf("line number must be positive")
	}

	return &FilePosition{
		FilePath: filePath,
		Line:     lineIdx,
	}, nil
}
