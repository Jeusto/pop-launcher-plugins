package ai

import (
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
)

type ChatAPI struct {
	Client    *openai.Client
	ApiKey    string
	MaxTokens int
	Prompt    string
	Request   openai.ChatCompletionRequest
}

func (c *ChatAPI) GetResponse(query string) (<-chan string, error) {
	c.Request.Messages = append(c.Request.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: query,
	})

	stream, err := c.createChatCompletionStream(c.Request)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)

	go func() {
		defer close(ch)
		full_result := ""

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			ch <- response.Choices[0].Delta.Content
			full_result += response.Choices[0].Delta.Content
		}

		c.Request.Messages = append(c.Request.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: full_result,
		})
	}()

	return ch, nil
}

func (c *ChatAPI) createChatCompletionStream(req openai.ChatCompletionRequest) (*openai.ChatCompletionStream, error) {
	return c.Client.CreateChatCompletionStream(context.Background(), req)
}

func New(apiKey string, maxTokens int, prompt string) *ChatAPI {
	return &ChatAPI{
		ApiKey:    apiKey,
		MaxTokens: maxTokens,
		Prompt:    prompt,
		Client:    openai.NewClient(apiKey),
		Request: openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: maxTokens,
			Stream:    true,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			},
		},
	}
}
