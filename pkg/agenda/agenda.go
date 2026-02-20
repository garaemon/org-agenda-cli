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
	// Truncate to day for comparison, forcing UTC to avoid timezone issues
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

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

// AdjustDate aligns the start date based on the range type.
// For "week", it returns the preceding Sunday.
// For "month", it returns the first day of the month.
// For other ranges, it returns the date as is.
func AdjustDate(date time.Time, rangeType string) time.Time {
	// Normalize to midnight
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	switch rangeType {
	case "week":
		// Weekday(): Sunday is 0, ... Saturday is 6.
		// Subtract the weekday value to get to Sunday.
		offset := int(d.Weekday())
		return d.AddDate(0, 0, -offset)
	case "month":
		return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
	default:
		return d
	}
}

func SortItems(items []*item.Item, criteria string, desc bool) {
	sort.SliceStable(items, func(i, j int) bool {
		var less bool
		switch criteria {
		case "priority":
			// Priority: A > B > C > ""
			// We want to sort A, B, C, "" in ascending order of importance?
			// Usually "sort" implies ascending rank.
			// But for priority, A is "smaller" char but "higher" priority.
			// Let's define ascending as A -> B -> C -> None.
			p1 := items[i].Priority
			p2 := items[j].Priority
			if p1 == "" && p2 != "" {
				less = false // "" is > "A" (in terms of being last)
			} else if p1 != "" && p2 == "" {
				less = true // "A" is < ""
			} else {
				less = p1 < p2
			}
		case "date":
			// Date: earliest deadline/schedule first.
			// We prioritize Deadline over Scheduled.
			t1 := getDate(items[i])
			t2 := getDate(items[j])
			if t1 == nil && t2 != nil {
				less = false
			} else if t1 != nil && t2 == nil {
				less = true
			} else if t1 == nil && t2 == nil {
				less = items[i].LineNumber < items[j].LineNumber // fallback
			} else {
				less = t1.Before(*t2)
			}
		case "status":
			// Status: TODO > WAITING > DONE > ""
			// Map status to integer
			s1 := getStatusRank(items[i].Status)
			s2 := getStatusRank(items[j].Status)
			less = s1 < s2
		default:
			// Default to file order (roughly)
			// But since we might merge multiple files, maybe just keep existing order?
			// Stable sort preserves order if we return false?
			// But let's verify stable sort behavior.
			// If we return item[i].LineNumber < item[j].LineNumber, it sorts by line.
			// But if files are different?
			return i < j // keep original order
		}

		if desc {
			return !less
		}
		return less
	})
}

func getDate(it *item.Item) *time.Time {
	if it.Deadline != nil {
		return it.Deadline
	}
	return it.Scheduled
}

func getStatusRank(status string) int {
	switch status {
	case item.StatusTodo:
		return 0
	case item.StatusWaiting:
		return 1
	case item.StatusDone:
		return 2
	default:
		return 3
	}
}
