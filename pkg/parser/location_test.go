package parser

import (
	"testing"
)

func TestParseFilePosition(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    *FilePosition
		wantErr bool
	}{
		{
			name: "Valid input",
			arg:  "todo.org:10",
			want: &FilePosition{
				FilePath: "todo.org",
				Line:     10,
			},
			wantErr: false,
		},
		{
			name:    "Invalid format (no colon)",
			arg:     "todo.org",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid format (multiple colons)",
			arg:     "path/to/file:10:20",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid line number (not integer)",
			arg:     "todo.org:abc",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid line number (zero)",
			arg:     "todo.org:0",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid line number (negative)",
			arg:     "todo.org:-1",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty file path",
			arg:     ":10",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFilePosition(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFilePosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.FilePath != tt.want.FilePath {
					t.Errorf("ParseFilePosition() FilePath = %v, want %v", got.FilePath, tt.want.FilePath)
				}
				if got.Line != tt.want.Line {
					t.Errorf("ParseFilePosition() Line = %v, want %v", got.Line, tt.want.Line)
				}
			}
		})
	}
}
