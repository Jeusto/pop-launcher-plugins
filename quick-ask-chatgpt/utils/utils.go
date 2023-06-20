package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const CONFIG_FILE = "config.json"
const PLUGIN_PATH = "/.local/share/pop-launcher/plugins/quick-ask-chatgpt/"

func WrapText(s string, line_width int) string {
	words := strings.Fields(s)

	// strings.Fields() removes trailing spaces, add one back if necessary
	if strings.HasSuffix(s, " ") {
		words = append(words, "")
	}

	if len(words) == 0 {
		return ""
	}

	wrapped := words[0]
	spaceLeft := line_width - len(wrapped)

	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = line_width - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}

	return wrapped
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
