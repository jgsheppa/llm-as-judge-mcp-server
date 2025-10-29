package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaClient struct {
	*BaseClient[*http.Client]
	baseURL string
}

func NewOllamaClient(apiKey, model, prompt string) LLMClient {
	return &OllamaClient{
		BaseClient: NewBaseClient(apiKey, model, prompt, func(key string) (*http.Client, error) {
			return &http.Client{}, nil
		}),
		baseURL: "http://localhost:11434",
	}
}

type OllamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaGenerateResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func (o *OllamaClient) Judge(ctx context.Context, question, response, evaluationFocus string) (string, error) {
	if err := o.HasError(); err != nil {
		return "", fmt.Errorf("ollama client initialization error: %w", err)
	}

	prompt := o.GetPrompt()

	reqBody := OllamaGenerateRequest{
		Model:  o.Model,
		Prompt: fmt.Sprintf(prompt, question, response, evaluationFocus),
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := o.GetClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp OllamaGenerateResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w. Body: %s", err, string(body))
	}

	if ollamaResp.Response == "" {
		return "", fmt.Errorf("ollama returned empty response. Full response: %+v", ollamaResp)
	}

	return ollamaResp.Response, nil
}
