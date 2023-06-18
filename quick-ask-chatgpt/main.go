package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"quick-ask-chatgpt/utils"

	"github.com/sashabaranov/go-openai"
)

type App struct {
	response string
}

const PROMPT = "Respond in a concise way to the following user query. The result should be in a single line with no special formatting: "
const CHAR_LIMIT = 80

type PluginSearchResult struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Keywords    interface{} `json:"keywords"`
	Icon        interface{} `json:"icon"`
	Exec        interface{} `json:"exec"`
	Window      interface{} `json:"window"`
}

type Append struct {
	Result PluginSearchResult `json:"Append"`
}

func (a *App) activate(index int) {
	utils.CopyToClipboard(a.response)
	fmt.Println("\"Close\"")
}

func (a *App) search(query string) {
	initialOutput := Append{
		Result: PluginSearchResult{
			Name:        "ChatGPT says :",
			Description: "",
		},
	}

	initialOutputBytes, _ := json.Marshal(initialOutput)

	fmt.Println("\"Clear\"")
	fmt.Printf("%s\n", initialOutputBytes)
	fmt.Println("\"Finished\"")

	api_key, _ := utils.RetrieveApiKey()
	c := openai.NewClient(api_key)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 500,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: PROMPT + query,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		a.response += response.Choices[0].Delta.Content

		// cmd := exec.Command("notify-send", a.response)
		// cmd.Run()

		output := Append{
			Result: PluginSearchResult{
				Name:        "ChatGPT says : (hit enter to copy)",
				Description: utils.SplitLongString(a.response, CHAR_LIMIT),
			},
		}

		outputBytes, _ := json.Marshal(output)

		fmt.Println("\"Clear\"")
		fmt.Printf("%s\n", outputBytes)
		fmt.Println("\"Finished\"")
	}

}

func main() {
	app := App{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		var request map[string]interface{}
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			continue
		}

		if _, ok := request["Search"]; ok {
			app.search(request["Search"].(string))
		} else if _, ok := request["Activate"]; ok {
			app.activate(int(request["Activate"].(float64)))
		}
	}
}
