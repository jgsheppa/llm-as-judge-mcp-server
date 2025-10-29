package client

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicClient struct {
	*BaseClient[anthropic.Client]
}

func NewAnthropicClient(apiKey, model, prompt string) LLMClient {
	return &AnthropicClient{
		BaseClient: NewBaseClient(apiKey, model, prompt, func(key string) (anthropic.Client, error) {
			return anthropic.NewClient(option.WithAPIKey(key)), nil
		}),
	}
}

func (a *AnthropicClient) Judge(ctx context.Context, question, response, evaluationFocus string) (string, error) {
	client := a.GetClient()
	prompt := a.GetPrompt()

	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(fmt.Sprintf(prompt, question, response, evaluationFocus))),
		},
		Model: anthropic.Model(a.Model),
	})
	if err != nil {
		return "", err
	}

	var result string
	for _, c := range message.Content {
		result += c.Text
	}

	return result, nil
}
