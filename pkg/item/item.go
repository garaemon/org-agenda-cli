package item

import "time"

const (
	StatusTodo    = "TODO"
	StatusDone    = "DONE"
	StatusWaiting = "WAITING"
)

// Item represents an entry in an Org file.
type Item struct {
	Title      string
	Level      int
	Status     string // "TODO", "DONE", "WAITING", etc.
	Tags       []string
	Scheduled  *time.Time
	Deadline   *time.Time
	FilePath   string
	LineNumber int
	RawContent string // Body content
}
