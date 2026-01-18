package capture

import (
	"regexp"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name      string
		template  string
		content   string
		want      string // for exact matches
		wantRegex string // for matches with timestamp
	}{
		{
			name:     "Content only",
			template: "* %c",
			content:  "Hello",
			want:     "* Hello",
		},
		{
			name:     "No placeholders",
			template: "* Static",
			content:  "Ignored",
			want:     "* Static",
		},
		{
			name:      "Timestamp only",
			template:  "Time: %t",
			content:   "Ignored",
			wantRegex: `Time: \[\d{4}-\d{2}-\d{2} \w{3} \d{2}:\d{2}\]`,
		},
		{
			name:      "Both",
			template:  "* %t %c",
			content:   "My Note",
			wantRegex: `\* \[\d{4}-\d{2}-\d{2} \w{3} \d{2}:\d{2}\] My Note`,
		},
		{
			name:      "Extended Date Formats",
			template:  "%Y-%m-%d %H:%M:%S %A %a",
			content:   "Ignored",
			wantRegex: `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \w+ \w{3}`,
		},
		{
			name:      "Year and content",
			template:  "%Y: %c",
			content:   "Diary",
			wantRegex: `\d{4}: Diary`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.template, tt.content)
			if tt.want != "" {
				if got != tt.want {
					t.Errorf("Format() = %v, want %v", got, tt.want)
				}
			}
			if tt.wantRegex != "" {
				matched, err := regexp.MatchString(tt.wantRegex, got)
				if err != nil {
					t.Errorf("Regex error: %v", err)
				}
				if !matched {
					t.Errorf("Format() = %v, want regex %v", got, tt.wantRegex)
				}
			}
		})
	}
}
