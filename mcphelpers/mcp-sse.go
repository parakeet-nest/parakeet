package mcphelpers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/tools"
)

func GetMCPSSEClient(ctx context.Context, baseURL string) (*client.SSEMCPClient, *mcp.InitializeResult, error) {

	mcpClient, err := client.NewSSEMCPClient(baseURL + "/sse")
	//fmt.Println("🔴 mcpClient:", mcpClient, baseURL+"/sse")
	if err != nil {
		return nil, nil, &MCPSSEClientCreationError{Message: fmt.Sprintf("Failed to create client: %v", err)}
	}
	defer mcpClient.Close()

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	// Start the client
	if err := mcpClient.Start(ctx); err != nil {
		return nil, nil, &MCPSSEClientStartError{Message: fmt.Sprintf("Failed to start client: %v", err)}
	}

	//fmt.Println("🟢 mcpClient:", mcpClient)
	// Initialize
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "parakeet mcp client",
		Version: "1.0.0",
	}

	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		return nil, nil, &MCPSSEClientInitializationError{Message: fmt.Sprintf("Failed to initialize client: %v", err)}
	}

	return mcpClient, initResult, nil
}

func GetSSETools(ctx context.Context, mcpClient *client.SSEMCPClient) ([]llm.Tool, error) {
	toolsRequest := mcp.ListToolsRequest{}
	mcpTools, err := mcpClient.ListTools(ctx, toolsRequest)

	if err != nil {
		return nil, &MCPGetToolsError{Message: fmt.Sprintf("Failed to list tools: %v", err)}
	}

	// Convert mcp.Tools to llm.Tools
	// TODO: Handle errors during conversion
	ollamaTools := tools.ConvertMCPTools(mcpTools.Tools)
	return ollamaTools, nil
}

func CallSSETool(ctx context.Context, mcpClient *client.SSEMCPClient, functionName string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
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
