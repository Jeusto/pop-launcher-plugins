package ai

import (
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
)

type ChatAPI struct {
	ApiKey    string
	MaxTokens int
	Prompt    string
}

func (c *ChatAPI) GetResponse(query string) (<-chan string, error) {
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: c.MaxTokens,
		Stream:    true,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: c.Prompt + query,
			},
		},
	}

	stream, err := c.createChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)

	go func() {
		defer close(ch)
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			ch <- response.Choices[0].Delta.Content
		}
	}()

	return ch, nil
}

func (c *ChatAPI) createChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionStream, error) {
	client := openai.NewClient(c.ApiKey)
	return client.CreateChatCompletionStream(ctx, req)
}

func New(apiKey string, maxTokens int, prompt string) *ChatAPI {
	return &ChatAPI{
		ApiKey:    apiKey,
		MaxTokens: maxTokens,
		Prompt:    prompt,
	}
}
