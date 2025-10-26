package client

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/jgsheppa/llm-as-judge-mcp-server/internal/prompts"
)

type AnthropicClient struct {
	*BaseClient[anthropic.Client]
}

func NewAnthropicClient(apiKey, model string) LLMClient {
	return &AnthropicClient{
		BaseClient: NewBaseClient(apiKey, model, func(key string) (anthropic.Client, error) {
			return anthropic.NewClient(option.WithAPIKey(key)), nil
		}),
	}
}

func (a *AnthropicClient) Judge(ctx context.Context, question, response, evaluationFocus string) (string, error) {
	client := a.GetClient()

	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(fmt.Sprintf(prompts.JudgePrompt, question, response, evaluationFocus))),
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
