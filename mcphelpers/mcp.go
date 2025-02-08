package mcphelpers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/tools"
)


// GetMCPClient creates and initializes a new MCP client.
//
// Parameters:
//   - ctx: The context for managing the lifecycle of the MCP client.
//   - command: The command to execute for the MCP client.
//   - env: The environment variables to set for the MCP client process.
//   - args: Additional arguments to pass to the MCP client command.
//
// Returns:
//   - *client.StdioMCPClient: The initialized MCP client.
//   - *mcp.InitializeResult: The result of the initialization process.
//   - error: An error if the client creation or initialization fails.
func GetMCPClient(ctx context.Context, command string, env []string, args ...string) (*client.StdioMCPClient, *mcp.InitializeResult, error) {
	// Create a new mcp client
	mcpClient, err := client.NewStdioMCPClient(command, env, args...)

	if err != nil {
		return nil, nil, &MCPClientCreationError{Message: fmt.Sprintf("Failed to create client: %v", err)}
	}

	// Initialize the client
	//fmt.Println("🚀 Initializing mcp client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "parakeet mcp client",
		Version: "1.0.0",
	}

	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		return nil, nil, &MCPClientInitializationError{Message: fmt.Sprintf("Failed to initialize client: %v", err)}
	}
	/*
		fmt.Printf(
			"🎉 Initialized with server: %s %s\n\n",
			initResult.ServerInfo.Name,
			initResult.ServerInfo.Version,
		)
	*/

	return mcpClient, initResult, nil
}

// GetTools retrieves a list of tools from the provided MCP client.
// It sends a request to list the tools and converts the received MCP tools
// to LLM tools format.
//
// Parameters:
//   - mcpClient: A pointer to an instance of client.StdioMCPClient.
//
// Returns:
//   - A slice of llm.Tool containing the converted tools.
//   - An error if there is any issue during the request or conversion process.
func GetTools(mcpClient *client.StdioMCPClient) ([]llm.Tool, error) {
	// List Tools
	//fmt.Println("🛠️ Available tools...")
	toolsRequest := mcp.ListToolsRequest{}
	mcpTools, err := mcpClient.ListTools(context.Background(), toolsRequest)
	if err != nil {
		return nil, &MCPGetToolsError{Message: fmt.Sprintf("Failed to list tools: %v", err)}
	}
	/*
		for _, tool := range mcpTools.Tools {
			fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
			fmt.Println("Arguments:", tool.InputSchema.Properties)
		}
	*/
	// Convert mcp.Tools to llm.Tools
	// TODO: Handle errors during conversion
	ollamaTools := tools.ConvertMCPTools(mcpTools.Tools)
	return ollamaTools, nil
}

// CallTool sends a request to the MCP server to call a specified tool function with given arguments.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - mcpClient: The client used to communicate with the MCP server.
//   - functionName: The name of the tool function to be called.
//   - arguments: A map containing the arguments to be passed to the tool function.
//
// Returns:
//   - *mcp.CallToolResult: The result of the tool function call.
//   - error: An error object if the call fails, otherwise nil.
func CallTool(ctx context.Context, mcpClient *client.StdioMCPClient, functionName string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	// 🖐️ Call the mcp server
	//fmt.Println("📣 calling", toolCall.Function.Name, "...")

	fetchRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	fetchRequest.Params.Name = functionName
	fetchRequest.Params.Arguments = arguments

	mcpResult, err := mcpClient.CallTool(ctx, fetchRequest)
	if err != nil {
		return nil, &MCPToolCallError{Message: fmt.Sprintf("Failed to call tool: %v", err)}
	}
	return mcpResult, nil
}

// GetTextFromResult extracts the text content from an mcp.CallToolResult object.
// It returns the extracted text as a string and an error if the extraction fails.
//
// Parameters:
// - mcpResult: A pointer to an mcp.CallToolResult object containing the result data.
//
// Returns:
// - string: The extracted text content.
// - error: An error if the extraction fails, which can occur in the following cases:
//   - The content slice is empty.
//   - The first element of the content slice is not a map.
//   - The map does not contain a "text" key.
//   - The value associated with the "text" key is not a string.
func GetTextFromResult(mcpResult *mcp.CallToolResult) (string, error) {
	//fmt.Println("📣 Result:", mcpResult.Result)
	//return mcpResult.Content[0].(map[string]interface{})["text"].(string)
	var finalText string
	if len(mcpResult.Content) == 0 {
		return "", &MCPResultExtractionError{Message: "content is empty"}
	}

	contentMap, ok := mcpResult.Content[0].(map[string]interface{})
	if !ok {
		return "", &MCPResultExtractionError{Message: "content[0] is not a map"}
	}

	textInterface, ok := contentMap["text"]
	if !ok {
		return "", &MCPResultExtractionError{Message: "no 'text' key in map"}
	}

	finalText, ok = textInterface.(string)
	if !ok {
		return "", &MCPResultExtractionError{Message: "text is not a string"}
	}
	return finalText, nil

}
