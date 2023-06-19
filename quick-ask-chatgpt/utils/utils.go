package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.design/x/clipboard"
	"os"
	"strings"
	"time"
	"unicode"
)

const CONFIG_FILE = "config.json"
const PLUGIN_PATH = "/.local/share/pop-launcher/plugins/quick-ask-chatgpt/"

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

func RetrieveApiKey() (string, error) {
	home := os.Getenv("HOME")
	data, read_file_error := os.ReadFile(home + PLUGIN_PATH + CONFIG_FILE)

	var api_key string = ""
	var config map[string]interface{}
	parse_error := json.Unmarshal(data, &config)

	if read_file_error != nil || parse_error != nil {
		return "", errors.New("error reading or parsing config file")
	}

	if config["OPENAI_API_KEY"] == nil {
		return "", errors.New("OPENAI_API_KEY not found in config file")
	}

	api_key = config["OPENAI_API_KEY"].(string)
	if api_key == "" {
		return "", errors.New("OPENAI_API_KEY in the config is empty")
	}

	return api_key, nil
}

func LogToFile(file string, s string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	f.WriteString("[" + timestamp + "] " + s + "\n")
}

func Ternary[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
}
