package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"mcp-curl",
		"1.0.0",
	)

	// Add a tool
	tool := mcp.NewTool("use_curl",
		mcp.WithDescription("fetch this webpage"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("url of the webpage to fetch"),
		),
	)

	// Add a tool handler
	s.AddTool(tool, curlHandler)

	fmt.Println("🚀 Server started")
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("😡 Server error: %v\n", err)
	}
	fmt.Println("👋 Server stopped")
}

func curlHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return mcp.NewToolResultError("url must be a string"), nil
	}
	cmd := exec.Command("curl", "-s", url)
	output, err := cmd.Output()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	content := string(output)

	return mcp.NewToolResultText(content), nil
}
