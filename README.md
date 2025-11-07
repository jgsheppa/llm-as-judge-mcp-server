# llm-as-judge-mcp-server

The LLM as Judge MCP Server enables users of LLMs to get a second opinion for an LLM's response. The MCP tool sends the user's question, the LLM's response, and an optional focus for the evaluation to a second LLM for evaluation. This second opinion can be used to then improve an LLM's response and give users another perspective regarding the original LLM's response.

## Installation

You can run this MCP server using `node` with the following command. This can also be used in MCP configuration files, which can be seen in the setup section below.

```
npx -y @jgsheppa/llm-as-judge-mcp-server
```

You can also download a binary which matches your machine's architecture to start using `llm-as-judge-mcp-server`.

## Setup

To set up the `llm-as-judge-mcp-server`, you can define the MCP server in a JSON file wherever your LLM has access to MCP servers.

```json
{
  "mcpServers": {
    "llm-as-judge": {
      "command": "npx",
      "args": [
        "-y",
        "@jgsheppa/llm-as-judge-mcp-server",
        "stdio",
        "-p", 
        "gemini", 
        ],
      "env": {
        "GEMINI_API_KEY": "your-api-key",
      }
    },
  }
}

```

To customize your model and prompt, you configuration would look like this:

```json
{
  "mcpServers": {
    "llm-as-judge": {
      "command": "npx",
      "args": [
        "-y",
        "@jgsheppa/llm-as-judge-mcp-server",
        "stdio",
        "-p", 
        "gemini", 
        "-m", 
        "gemini-2.5-flash", 
        "--prompt-path",
        "/Users/firstlast/Desktop/PROMPT.md"
        ],
      "env": {
        "GEMINI_API_KEY": "your-api-key",
      }
    },
  }
}

```

To define your provider and model, you can pass them as arguments to the MCP server. And while there is a default prompt for the LLM as judge, you can also pass your own custom prompt by providing the full filepath to the `--prompt-path` argument.

## Features

### Providers

Currently three providers are available for this MCP server:

- Anthropic
- Gemini
- Ollama
- OpenAI

More providers can and will be added in the future.

### Models

Any models offered by the previously mentioned providers can be used as LLM judges. It is up to you to decide which model works best for your use-case, but a bigger, frontier model is not necessarily the best option to evaluate the response of another frontier model. Smaller, less expensive models can be great options as well.

#### Default Models

To try and improve the experience of using this MCP server, there are default models for each provider. The models were chosen because they are the most cost-efficient. These might change in the future, but since they are customizable, these seem like fine choices for now.

| Provider  | Default Model      |
|-----------|-------------------|
| Anthropic | claude-haiku-4-5  |
| Gemini    | gemini-2.5-flash  |
| Ollama    | gemma3:4b         |
| OpenAI    | gpt-5-mini        |

### Prompts

This MCP server offers a default prompt for an LLM to evaluate another LLM's response, but you can also provide your own prompt for further customization.