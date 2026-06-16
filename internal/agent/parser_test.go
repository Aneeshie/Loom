package agent

import "testing"

func TestParseLogLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		level   string
		message string
		wantErr bool
	}{
		{
			name:    "valid INFO log",
			input:   "[INFO] user logged in",
			level:   "INFO",
			message: "user logged in",
			wantErr: false,
		},
		{
			name:    "invalid no brackets",
			input:   "this has no brackets",
			level:   "",
			message: "",
			wantErr: true,
		},
		{
			name:    "message with trailing spaces",
			input:   "[DEBUG]   lots of spaces   ",
			level:   "DEBUG",
			message: "lots of spaces", // TrimSpace only trims leading/trailing, not middle
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, message, err := ParseLogLine(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLogLine() error = %v, wantErr %v", err, tt.wantErr)
			}

			if level != tt.level {
				t.Errorf("ParseLogLine() level = %v, want %v", level, tt.level)
			}

			if message != tt.message {
				t.Errorf("ParseLogLine() message = %v, want %v", message, tt.message)
			}
		})
	}
}
