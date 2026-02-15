package item

import "time"

const (
	StatusTodo    = "TODO"
	StatusDone    = "DONE"
	StatusWaiting = "WAITING"
)

// Item represents an entry in an Org file.
type Item struct {
	Title      string     `json:"title"`
	Level      int        `json:"level"`
	Status     string     `json:"status"`
	Priority   string     `json:"priority,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	Scheduled  *time.Time `json:"scheduled,omitempty"`
	Deadline   *time.Time `json:"deadline,omitempty"`
	FilePath   string     `json:"filePath"`
	LineNumber int        `json:"lineNumber"`
	RawContent string     `json:"rawContent,omitempty"`
}
