package capture

import (
	"strings"
	"time"
)

// Format replaces placeholders in the template with actual values.
// %c: content
// %t: timestamp [YYYY-MM-DD Mon HH:MM]
func Format(template string, content string) string {
	// Default Org timestamp format: [2006-01-02 Mon 15:04]
	now := time.Now()
	timestamp := now.Format("[2006-01-02 Mon 15:04]")

	replacer := strings.NewReplacer(
		"%c", content,
		"%t", timestamp,
	)

	return replacer.Replace(template)
}
