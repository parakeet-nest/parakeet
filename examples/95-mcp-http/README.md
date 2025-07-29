# MCP HTTP Client Example

This example demonstrates how to use the MCP HTTP client to connect to an MCP server via HTTP transport.

## Prerequisites

1. Start the MCP HTTP server (in `mcp-server/` directory):
   ```bash
   cd mcp-server
   go run main.go
   ```
   The server will run on port 9090 by default.

2. Make sure you have Ollama running with a model that supports function calling (e.g., qwen2.5:latest).

## Running the Example

1. Set up your environment variables in `.env`:
   ```env
   OLLAMA_HOST=http://localhost:11434
   LLM_WITH_TOOLS_SUPPORT=qwen2.5:latest
   MCP_HOST=http://localhost:9090
   ```

2. Run the client:
   ```bash
   go run main.go
   ```

## What it does

The example:
1. Connects to the MCP HTTP server
2. Initializes the client
3. Lists available tools from the server
4. Uses a chat completion with the LLM to demonstrate tool calling
5. When the LLM wants to call a tool, it executes the tool via the MCP client
6. Shows the result of the tool execution

The HTTP server provides a `rool_dices` tool that simulates rolling dice, which the LLM can call when asked to roll dice.