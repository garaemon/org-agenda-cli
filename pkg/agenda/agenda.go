package agenda

import (
	"sort"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/item"
)

// FilterItemsByRange returns items that have a schedule or deadline within the range [start, end].
func FilterItemsByRange(items []*item.Item, start, end time.Time) []*item.Item {
	var filtered []*item.Item
	for _, it := range items {
		if it.Scheduled != nil {
			if isWithin(it.Scheduled, start, end) {
				filtered = append(filtered, it)
				continue
			}
		}
		if it.Deadline != nil {
			if isWithin(it.Deadline, start, end) {
				filtered = append(filtered, it)
				continue
			}
		}
	}
	return filtered
}

func isWithin(t *time.Time, start, end time.Time) bool {
	// Truncate to day for comparison
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	return (d.After(s) || d.Equal(s)) && (d.Before(e) || d.Equal(e))
}

// ExtractUniqueTags returns a sorted list of unique tags from the given items.
func ExtractUniqueTags(items []*item.Item) []string {
	tagMap := make(map[string]bool)
	for _, it := range items {
		for _, tag := range it.Tags {
			if tag != "" {
				tagMap[tag] = true
			}
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}
