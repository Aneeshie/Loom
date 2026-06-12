package agent

import (
	"fmt"
	"strings"
)

func ParseLogLine(logLine string) (string, string, error) {

	start := strings.Index(logLine, "[")
	end := strings.Index(logLine, "]")

	if start != -1 && end != -1 && end > start {
		level := logLine[start+1 : end]

		message := strings.TrimSpace(logLine[end+1:])

		return level, message, nil

	}
	return "", "", fmt.Errorf("invalid format")
}
