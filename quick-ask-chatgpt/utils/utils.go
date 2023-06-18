package utils

import (
	"golang.design/x/clipboard"
	"strings"
	"unicode"
)

func CopyToClipboard(content string) (string, error) {
	err := clipboard.Init()

	if err != nil {
		return "", err
	}

	clipboard.Write(clipboard.FmtText, []byte(content))

	return "Successfully copied to clipboard", nil
}

func SplitLongString(s string, char_limit uint) string {
	var result strings.Builder
	var line strings.Builder
	var word strings.Builder
	var lineWidth int

	for _, r := range s {
		if unicode.IsSpace(r) {
			// Found a space or newline character.
			// Add the current word to the current line,
			// if it fits, otherwise start a new line.
			if lineWidth+len(word.String())+1 <= int(char_limit) {
				line.WriteString(word.String())
				line.WriteRune(' ')
				lineWidth += len(word.String()) + 1
			} else {
				result.WriteString(strings.TrimRight(line.String(), " "))
				result.WriteRune('\n')
				line.Reset()
				line.WriteString(word.String())
				line.WriteRune(' ')
				lineWidth = len(word.String()) + 1
			}
			word.Reset()
		} else {
			// Found a non-space character, add it to the current word.
			word.WriteRune(r)
		}
	}

	// Add the last word and line to the result.
	if lineWidth+len(word.String()) <= 80 {
		line.WriteString(word.String())
		result.WriteString(strings.TrimRight(line.String(), " "))
	} else {
		result.WriteString(strings.TrimRight(line.String(), " "))
		result.WriteRune('\n')
		result.WriteString(word.String())
	}

	return result.String()
}
