package capture

import (
	"strings"
	"time"
)

// Format replaces placeholders in the template with actual values.
// Supports:
// %c: content
// %t: timestamp [YYYY-MM-DD Mon HH:MM]
// %Y: Year (2006)
// %y: Year (06)
// %m: Month (01)
// %d: Day (02)
// %H: Hour (15)
// %M: Minute (04)
// %S: Second (05)
// %A: Day of week (Monday)
// %a: Day of week (Mon)
func Format(template string, content string) string {
	now := time.Now()

	replacements := []string{
		"%c", content,
		"%t", now.Format("<2006-01-02 Mon 15:04>"),
		"%Y", now.Format("2006"),
		"%y", now.Format("06"),
		"%m", now.Format("01"),
		"%d", now.Format("02"),
		"%H", now.Format("15"),
		"%M", now.Format("04"),
		"%S", now.Format("05"),
		"%A", now.Format("Monday"),
		"%a", now.Format("Mon"),
	}

	replacer := strings.NewReplacer(replacements...)
	return replacer.Replace(template)
}
