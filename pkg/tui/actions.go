package tui

import (
	"fmt"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/agenda"
	"github.com/garaemon/org-agenda-cli/pkg/item"
)

func ToggleStatus(it *item.Item) error {
	switch it.Status {
	case item.StatusTodo:
		it.Status = item.StatusDone
	case item.StatusDone:
		it.Status = item.StatusWaiting
	case item.StatusWaiting:
		it.Status = item.StatusTodo
	default:
		it.Status = item.StatusTodo
	}
	// SaveItem updates the file content based on the item state
	if err := agenda.SaveItem(it); err != nil {
		return err
	}
	return nil
}

func CyclePriority(it *item.Item) error {
	switch it.Priority {
	case "A":
		it.Priority = "B"
	case "B":
		it.Priority = "C"
	case "C":
		it.Priority = ""
	case "":
		it.Priority = "A"
	default:
		it.Priority = "A"
	}
	if err := agenda.SaveItem(it); err != nil {
		return err
	}
	return nil
}

func UpdateTimestamp(it *item.Item, key string, val string) error {
	// Validate format
	t, err := time.Parse("2006-01-02", val)
	if err != nil {
		return fmt.Errorf("invalid date: %v", err)
	}

	if key == "SCHEDULED" {
		it.Scheduled = &t
	}
	if key == "DEADLINE" {
		it.Deadline = &t
	}

	if err := agenda.UpdateTimestamp(it, key, t); err != nil {
		return err
	}
	return nil
}
