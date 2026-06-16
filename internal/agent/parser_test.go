package agent

import "testing"

func TestParseLogLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		level   string
		message string
	}{
		{
			name:    "valid INFO log",
			input:   "[INFO] user logged in",
			level:   "INFO",
			message: "user logged in",
		},
		{

			name:    "invalid no brackets becomes RAW",
			input:   "this has no brackets",
			level:   "RAW",
			message: "this has no brackets",

		},
		{
			name:    "message with trailing spaces",
			input:   "[DEBUG]   lots of spaces   ",
			level:   "DEBUG",
			message: "lots of spaces", // TrimSpace only trims leading/trailing, not middle
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, message := ParseLogLine(tt.input)

			if level != tt.level {
				t.Errorf("ParseLogLine() level = %v, want %v", level, tt.level)
			}

			if message != tt.message {
				t.Errorf("ParseLogLine() message = %v, want %v", message, tt.message)
			}
		})
	}
}
