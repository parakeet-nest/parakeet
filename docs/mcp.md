# Model Context Protocol

!!! info "🚧 work in progress"

Integration of `github.com/mark3labs/mcp-go/mcp` and `github.com/mark3labs/mcp-go/client`

## Helpers

- `mcphelpers.GetMCPClient(ctx context.Context, command string, env []string, args ...string) (*client.StdioMCPClient, *mcp.InitializeResult, error)`
- `mcphelpers.GetTools(mcpClient *client.StdioMCPClient) ([]llm.Tool, error)`
- `tools.ConvertMCPTools` to convert the MCP tools list to a list compliant with Ollama LLM tools. (used by `GetTools`)
- `mcphelpers.CallTool(ctx context.Context, mcpClient *client.StdioMCPClient, functionName string, arguments map[string]interface{}) (*mcp.CallToolResult, error)`
- `mcphelpers.GetTextFromResult(mcpResult *mcp.CallToolResult) (string, error)`


!!! note
	👀 you will find a complete example in:

    - [examples/67-mcp](https://github.com/parakeet-nest/parakeet/tree/main/examples/67-mcp)
