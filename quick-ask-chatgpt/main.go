package main

import (
	"quick-ask-chatgpt/ai"
	"quick-ask-chatgpt/pop"
	"quick-ask-chatgpt/utils"
)

const MAX_TOKENS = 1000
const CHAR_LIMIT = 80

const HELP_TEXT = "Hi! How can I help you?\nStart typing and hit enter to send your query to ChatGPT."
const PROMPT = `Provide a concise response to the following user question
in a single line without any special formatting. Focus on the most relevant
information. Respond like you're directly chatting with the person. Never try
to complete user's query, you can tell him if the question is not complete or
clear enough: `

type Plugin struct {
	api_response       string
	query              string
	chat_api           *ai.ChatAPI
	previous_questions []string
	previous_answers   []string
}

func (plugin *Plugin) activate(index int) {
	plugin.api_response = ""
	pop.ClearInput()

	responseCh, err := plugin.chat_api.GetResponse(plugin.query)
	if err != nil {
		pop.ShowErrorMessage(err.Error())
		return
	}

	for response := range responseCh {
		plugin.api_response += response
		plugin.api_response = utils.SplitLongString(plugin.api_response, CHAR_LIMIT)

		pop.ShowSingleResult(pop.PluginSearchResult{
			Name: plugin.api_response,
		})
	}
}

func (plugin *Plugin) search(query string) {
	plugin.query = query

	if len(plugin.api_response) > 0 {
		pop.ShowSingleResult(pop.PluginSearchResult{
			Name: plugin.api_response,
		})
	} else {
		pop.ShowSingleResult(pop.PluginSearchResult{
			Name: HELP_TEXT,
		})
	}
}

func main() {
	api_key, err := utils.RetrieveApiKey()
	if err != nil {
		pop.ShowErrorMessage(err.Error())
		return
	}

	chat_api := ai.ChatAPI{
		ApiKey:    api_key,
		Prompt:    PROMPT,
		MaxTokens: MAX_TOKENS,
	}

	plugin := Plugin{
		api_response: "",
		chat_api:     &chat_api,
	}

	requests := make(chan pop.Request)
	go pop.HandleRequests(requests)

	for request := range requests {
		if request.Type == "search" {
			plugin.search(request.Query)
		} else if request.Type == "activate" {
			plugin.activate(request.ID)
		}
	}
}
