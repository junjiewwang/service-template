package internal

import (
	"strings"
)

// TrimEmptyLines removes empty lines from the beginning and end of text
func TrimEmptyLines(text string) string {
	lines := strings.Split(text, "\n")

	// Trim from start
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	// Trim from end
	end := len(lines)
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}

	return strings.Join(lines[start:end], "\n")
}

// IndentText adds indentation to each line of text
func IndentText(text string, indent string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}

// EnsureTrailingNewline ensures text ends with a newline
func EnsureTrailingNewline(text string) string {
	if !strings.HasSuffix(text, "\n") {
		return text + "\n"
	}
	return text
}
