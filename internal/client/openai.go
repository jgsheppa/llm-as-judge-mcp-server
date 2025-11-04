package client

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIClient struct {
	*BaseClient[openai.Client]
}

func NewOpenAIClient(apiKey, model, prompt string) LLMClient {
	return &OpenAIClient{
		BaseClient: NewBaseClient(apiKey, model, prompt, func(key string) (openai.Client, error) {
			return openai.NewClient(
				option.WithAPIKey(apiKey),
			), nil
		}),
	}
}

func (o *OpenAIClient) Judge(ctx context.Context, question, response, evaluationFocus string) (string, error) {
	prompt := o.GetPrompt()

	chatCompletion, err := o.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(fmt.Sprintf(prompt, question, response, evaluationFocus)),
		},
		Model: o.Model,
	})
	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI API")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}
