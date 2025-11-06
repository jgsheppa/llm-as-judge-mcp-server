package main

import (
	"fmt"
	"os"

	"github.com/jgsheppa/llm-as-judge-mcp-server/internal/client"
	"github.com/jgsheppa/llm-as-judge-mcp-server/internal/config"
	"github.com/jgsheppa/llm-as-judge-mcp-server/internal/handler"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "llm-as-judge-mcp-server",
	Short: "A command line interface for llm-as-judge MCP Server",
	Long: `llm-as-judge-mcp-server is a command line interface 
	to start an MCP server that uses an LLM as a judge for 
	evaluating responses. Users can specify the LLM provider,
	model, and prompt path via command line flags.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var stdioCommand = &cobra.Command{
	Use:   "stdio",
	Short: "Start stdio server",
	Long: `Start the MCP server over stdio for use with AI 
	clients which support the MCP protocol.`,
	Example: "stdio --provider=gemini",
	Run: func(cmd *cobra.Command, _ []string) {
		provider, _ := cmd.Flags().GetString("provider")
		model, _ := cmd.Flags().GetString("model")
		promptPath, _ := cmd.Flags().GetString("prompt-path")

		args := []string{provider, model, promptPath}

		cfg, err := config.Load(provider)
		if err != nil {
			fmt.Println("Configuration error:")
			fmt.Println(err)
			os.Exit(1)
		}

		apiKey := cfg.ProviderAPIKey

		llmClient := client.GetClientProvider(apiKey, args)

		judgeHandler := handler.NewJudgeHandler(llmClient)

		s := server.NewMCPServer(
			"llm-as-judge",
			"0.8.1",
			server.WithToolCapabilities(false),
		)

		s.AddTool(handler.NewTool(), judgeHandler.Handle)

		if err := server.ServeStdio(s); err != nil {
			fmt.Println("Server error:")
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(stdioCommand)

	stdioCommand.Flags().StringP("provider", "p", "", "The LLM provider to use as a judge (anthropic, openai, gemini, ollama)")
	stdioCommand.MarkFlagRequired("provider")
	stdioCommand.Flags().StringP("model", "m", "", "The model to use which is offered by the given provider")
	stdioCommand.Flags().String("prompt-path", "", "An optional path to your prompt")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
