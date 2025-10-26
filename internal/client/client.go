package client

import (
	"context"
)

type LLMClient interface {
	Judge(ctx context.Context, question, response, evaluationFocus string) (string, error)
}

type BaseClient[T any] struct {
	APIKey string
	Model  string
	client T
	err    error
}

func NewBaseClient[T any](apiKey string, model string, clientFactory func(string) (T, error)) *BaseClient[T] {
	client, err := clientFactory(apiKey)

	return &BaseClient[T]{
		APIKey: apiKey,
		Model:  model,
		client: client,
		err:    err,
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
