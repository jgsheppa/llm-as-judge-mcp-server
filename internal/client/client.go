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
