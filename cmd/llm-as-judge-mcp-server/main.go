package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jgsheppa/llm-as-judge-mcp-server/internal/client"
	"github.com/jgsheppa/llm-as-judge-mcp-server/internal/config"
	"github.com/jgsheppa/llm-as-judge-mcp-server/internal/handler"
	"github.com/mark3labs/mcp-go/server"
)

func GetDefaultProviderModel(provider string) string {
	var model string
	switch provider {
	case "anthropic":
		model = "claude-haiku-4-5"
	case "gemini":
		model = "gemini-2.5-flash"
	case "ollama":
		model = "gemma3:4b"
	default:
		model = ""
	}
	return model
}

func main() {
	provider := flag.String("provider", "gemini", "the LLM provider to use as a judge (anthropic, openai, gemini)")
	defaultModel := GetDefaultProviderModel(*provider)
	model := flag.String("model", defaultModel, "the model for the given provider")
	promptPath := flag.String("prompt", "", "an optional path to your prompt")

	flag.Parse()

	cfg, err := config.Load(*provider)
	if err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		os.Exit(1)
	}

	apiKey := cfg.ProviderAPIKey

	var llmClient client.LLMClient
	switch *provider {
	case "anthropic":
		llmClient = client.NewAnthropicClient(apiKey, *model, *promptPath)
	case "gemini":
		llmClient = client.NewGeminiClient(apiKey, *model, *promptPath)
	case "ollama":
		llmClient = client.NewOllamaClient(apiKey, *model, *promptPath)
	}

	judgeHandler := handler.NewJudgeHandler(llmClient)

	s := server.NewMCPServer(
		"llm-as-judge",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	s.AddTool(handler.NewTool(), judgeHandler.Handle)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
