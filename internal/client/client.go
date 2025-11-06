package client

import (
	"context"
	"os"

	"github.com/jgsheppa/llm-as-judge-mcp-server/internal/prompts"
)

type LLMClient interface {
	Judge(ctx context.Context, question, response, evaluationFocus string) (string, error)
}

type BaseClient[T any] struct {
	APIKey     string
	Model      string
	PromptPath string
	client     T
	err        error
}

func NewBaseClient[T any](apiKey string, model string, promptPath string, clientFactory func(string) (T, error)) *BaseClient[T] {
	client, err := clientFactory(apiKey)

	return &BaseClient[T]{
		APIKey:     apiKey,
		Model:      model,
		PromptPath: promptPath,
		client:     client,
		err:        err,
	}
}

func GetDefaultProviderModel(provider, model string) string {
	if model != "" {
		return model
	}

	switch provider {
	case "anthropic":
		model = "claude-haiku-4-5"
	case "gemini":
		model = "gemini-2.5-flash"
	case "ollama":
		model = "gemma3:4b"
	case "openai":
		model = "gpt-5-mini"
	default:
		model = ""
	}
	return model
}

func GetClientProvider(apiKey string, args []string) LLMClient {
	provider := args[0]
	var model, promptPath string
	if len(args) > 1 {
		model = args[1]
	}
	if len(args) > 2 {
		promptPath = args[2]
	}
	providerModel := GetDefaultProviderModel(provider, model)

	var llmClient LLMClient
	switch provider {
	case "anthropic":
		llmClient = NewAnthropicClient(apiKey, providerModel, promptPath)
	case "gemini":
		llmClient = NewGeminiClient(apiKey, providerModel, promptPath)
	case "ollama":
		llmClient = NewOllamaClient(apiKey, providerModel, promptPath)
	case "openai":
		llmClient = NewOpenAIClient(apiKey, providerModel, promptPath)
	}
	return llmClient
}

func (b *BaseClient[T]) GetClient() T {
	return b.client
}

func (b *BaseClient[T]) HasError() error {
	if b.err != nil {
		return b.err
	}
	return nil
}

func (b *BaseClient[T]) GetPrompt() string {
	if b.PromptPath == "" {
		return prompts.JudgePrompt

	}

	content, err := os.ReadFile(b.PromptPath)
	if err != nil {
		b.err = err
		return ""
	}

	return string(content)
}
