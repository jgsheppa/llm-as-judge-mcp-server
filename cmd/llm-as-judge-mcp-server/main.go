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

func main() {
	provider := flag.String("provider", "gemini", "the LLM provider to use as a judge (anthropic, openai, gemini)")
	defaultModel := client.GetDefaultProviderModel(*provider)
	model := flag.String("model", defaultModel, "the model for the given provider")
	promptPath := flag.String("prompt", "", "an optional path to your prompt")

	flag.Parse()

	cfg, err := config.Load(*provider)
	if err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		os.Exit(1)
	}

	apiKey := cfg.ProviderAPIKey

	llmClient := client.GetClientProvider(*provider, apiKey, *model, *promptPath)

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
