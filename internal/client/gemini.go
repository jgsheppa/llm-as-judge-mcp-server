package client

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiClient struct {
	*BaseClient[genai.Client]
}

func NewGeminiClient(apiKey, model, prompt string) LLMClient {
	return &GeminiClient{
		BaseClient: NewBaseClient(apiKey, model, prompt, func(key string) (genai.Client, error) {
			client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
				APIKey:  apiKey,
				Backend: genai.BackendGeminiAPI,
			})
			if err != nil {
				return genai.Client{}, err
			}
			return *client, nil
		}),
	}
}

func (g *GeminiClient) Judge(ctx context.Context, question, response, evaluationFocus string) (string, error) {
	var config *genai.GenerateContentConfig = &genai.GenerateContentConfig{Temperature: genai.Ptr[float32](0.5)}

	chat, err := g.client.Chats.Create(ctx, g.Model, config, nil)
	if err != nil {
		return "", err
	}

	prompt := g.GetPrompt()

	result, err := chat.SendMessage(ctx, genai.Part{Text: fmt.Sprintf(prompt, question, response, evaluationFocus)})
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
