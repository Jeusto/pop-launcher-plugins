package main

import (
	"bufio"
	"encoding/json"
	"os"

	"quick-ask-chatgpt/ai"
	"quick-ask-chatgpt/pop"
	"quick-ask-chatgpt/utils"

	"github.com/sashabaranov/go-openai"
)

const AI_MODEL = openai.GPT3Dot5Turbo
const PROMPT = "Respond in a concise way to the following user query. The result should be in a single line with no special formatting: "
const RESPONSE_TOKEN_LIMIT = 500
const RESPONSE_TITLE = "ChatGPT says :"
const CHAR_LIMIT = 80

type Plugin struct {
	api_response string
	chat_api     *ai.ChatAPI
}

func (plugin *Plugin) activate(index int) {
	utils.CopyToClipboard(plugin.api_response)
	pop.CloseLauncher()
}

func (plugin *Plugin) search(query string) {
	pop.ClearSearchResults()
	pop.AppendResult(pop.PluginSearchResult{
		Name:        RESPONSE_TITLE,
		Description: "",
	})
	pop.Finish()

	responseCh, err := plugin.chat_api.GetResponse(query)
	if err != nil {
		pop.ShowErrorMessage(err.Error())
		return
	}

	for response := range responseCh {
		plugin.api_response += response
		pop.ClearSearchResults()
		pop.AppendResult(pop.PluginSearchResult{
			Name:        RESPONSE_TITLE,
			Description: utils.SplitLongString(plugin.api_response, CHAR_LIMIT),
		})
		pop.Finish()
	}

}

func main() {
	api_key, err := utils.RetrieveApiKey()
	if err != nil {
		pop.ShowErrorMessage(err.Error())
		return
	}

	chat_api := ai.ChatAPI{
		ApiKey: api_key,
	}

	plugin := Plugin{
		api_response: "",
		chat_api:     &chat_api,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		var request map[string]interface{}
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			continue
		}

		if _, ok := request["Search"]; ok {
			plugin.search(request["Search"].(string))
		} else if _, ok := request["Activate"]; ok {
			plugin.activate(int(request["Activate"].(float64)))
		}
	}
}
