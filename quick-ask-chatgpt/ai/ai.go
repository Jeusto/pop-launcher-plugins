package ai

import (
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
)

const AI_MODEL = openai.GPT3Dot5Turbo
const PROMPT = "Respond in a concise way to the following user query. The result should be in a single line with no special formatting:"
const RESPONSE_TOKEN_LIMIT = 500
const CHAR_LIMIT = 80

type ChatAPI struct {
	ApiKey string
}

func (c *ChatAPI) GetResponse(query string) (<-chan string, error) {
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:     AI_MODEL,
		MaxTokens: RESPONSE_TOKEN_LIMIT,
		Stream:    true,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: PROMPT + query,
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

func New(apiKey string) *ChatAPI {
	return &ChatAPI{ApiKey: apiKey}
}
