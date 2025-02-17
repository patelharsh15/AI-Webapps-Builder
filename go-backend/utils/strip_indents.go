// utils/strip_indents.go
package utils

import (
	"strings"
)

func StripIndents(value string) string {
	lines := strings.Split(value, "\n")

	// Trim each line
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	// Join lines and trim start
	result := strings.TrimSpace(strings.Join(lines, "\n"))

	// Replace trailing newlines
	result = strings.TrimRight(result, "\r\n")

	return result
}
