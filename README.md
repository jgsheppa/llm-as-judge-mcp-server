# llm-as-judge-mcp-server

The LLM as Judge MCP Server enables users of LLMs to get a second opinion for an LLM's response. The MCP tool sends the user's question, the LLM's response, and an optional focus for the evaluation to a second LLM for evaluation. This second opinion can be used to then improve an LLM's response and give users another perspective regarding the original LLM's response.

## Installation

You can download a binary which matches your machine's architecture to start using `llm-as-judge-mcp-server`.

## Setup

To set up the `llm-as-judge-mcp-server`, you can define the MCP server in a JSON file wherever your LLM has access to MCP servers.

```json
{
  "mcpServers": {
    "llm-as-judge": {
      "command": "llm-as-judge-mcp-server",
      "args": [
        // User-defined provider
        "-provider", 
        "gemini", 
        // User-defined model
        "-model", 
        "gemini-2.5-flash", 
        // Custom prompt with filepath for Mac users
        "-prompt",
        "/Users/firstlast/Desktop/PROMPT.md"
        ],
      "env": {
        "GEMINI_API_KEY": "your-api-key",
      }
    },
  }
}

```

To define your provider and model, you can pass them as arguments to the MCP server. And while there is a default prompt for the LLM as judge, you can also pass your own custom prompt by providing the full filepath to the `-prompt` argument.

## Features

### Providers

Currently three providers are available for this MCP server:

- Gemini
- Anthropic
- Ollama

More providers can and will be added in the future.

### Models

Any models offered by the previously mentioned providers can be used as LLM judges. It is up to you to decide which model works best for your use-case, but a bigger, frontier model is not necessarily the best option to evaluate the response of another frontier model. Smaller, less expensive models can be great options as well.

### Prompts

This MCP server offers a default prompt for an LLM to evaluate another LLM's response, but you can also provide your own prompt for further customization.