package main

import (
	"context"
	"fmt"
	"log"
	"time"

	mcpsse "github.com/parakeet-nest/parakeet/mcp-sse"
)

func main() {

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new mcp client
	mcpClient, err := mcpsse.NewClient(ctx, "http://0.0.0.0:5001")
	defer mcpClient.Close()

	if err != nil {
		log.Fatalln("😡 error when creating the MCP client:", err)
	}

	// Start and initialize the client
	err = mcpClient.Start()

	if err != nil {
		log.Fatalln("😡 error when starting the MCP client:", err)
	}

	result, err := mcpClient.Initialize()

	if err != nil {
		log.Fatalln("😡 error when initializing the MCP client:", err)
	}

	fmt.Println("🚀 Initialized with server:", result.ServerInfo.Name, result.ServerInfo.Version)

	// ------------------------------
	//  List and read the ressources
	// ------------------------------

	// Get the list of resources
	resources, err := mcpClient.ListResources()
	if err != nil {
		log.Fatalln("😡", err)
	}

	// Print the list of available resources
	fmt.Println("🌍 Available Static Resources:")
	for _, resource := range resources {
		fmt.Printf("- Name: %s, URI: %s \n", resource.Name, resource.URI)
	}

	/*
		    🌍 Available Static Resources:
			- Name: tools system instructions, URI: tools-system://instructions, MIME Type:
			- Name: chat system instructions, URI: chat-system://instructions, MIME Type:
	*/

	fmt.Println()
	fmt.Println("📝 Resources content:")

	resourceResult, err := mcpClient.ReadResource("tools-system://instructions")
	if err != nil {
		log.Fatalln("😡 Failed to read resource:", err)
	}
	toolsSystemInstructions := resourceResult.Contents[0]["text"].(string)

	resourceResult, err = mcpClient.ReadResource("chat-system://instructions")
	if err != nil {
		log.Fatalln("😡", err)
	}
	chatSystemInstructions := resourceResult.Contents[0]["text"].(string)

	fmt.Println("- Tools System Instructions:", toolsSystemInstructions)
	fmt.Println("- Chat System Instructions:", chatSystemInstructions)

	/*
		📝 Resources content:
		- Tools System Instructions: You are a useful AI agent.
		Your job is to understand the user prompt ans decide if you need to use a tool to run external commands.
		Ignore all things not related to the usage of a tool

		- Chat System Instructions: You are a useful AI agent. your job is to answer the user prompt.
		If you detect that the user prompt is related to a tool, ignore this part and focus on the other parts.
	*/

	// ------------------------------
	//  List and read the prompts
	// ------------------------------

	// Get the list of prompts
	prompts, err := mcpClient.ListPrompts()
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println()
	fmt.Println("📣 Get the list of the prompts")
	for _, prompt := range prompts {
		fmt.Println("- Name:", prompt.Name, "Arguments:", prompt.Arguments)
	}

	fmt.Println()
	fmt.Println("📝 Fill the fetch-page prompt")

	fetchPrompt, err := mcpClient.GetAndFillPrompt("fetch-page", map[string]string{"url": "https://docker.com"})
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println(
		"📣 Filled Prompt:",
		"role:", fetchPrompt.Messages[0].Role,
		"content:", fetchPrompt.Messages[0].Content,
	)

	fmt.Println()
	fmt.Println("📝 Fill the summarize prompt")

	summarizePrompt, err := mcpClient.GetAndFillPrompt("summarize", map[string]string{"content": "[this is the content of the page]]"})
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println(
		"📣 Filled Prompt:",
		"role:", summarizePrompt.Messages[0].Role,
		"content:", summarizePrompt.Messages[0].Content,
	)

	// ------------------------------
	//  List and read the tools
	// ------------------------------
	fmt.Println()
	fmt.Println("🛠️ Get tools list from the MCP server")
	ollamaTools, err := mcpClient.ListTools()
	if err != nil {
		log.Fatalln("😡", err)
	}

	for _, tool := range ollamaTools {
		fmt.Println("🛠️ Tool:", tool.Function.Name)
		fmt.Println("  - Arguments:")
		for name, prop := range tool.Function.Parameters.Properties {
			fmt.Println("    - name", name, ":", prop.Type)
		}
	}

	fmt.Println()
	fmt.Println("🛠️ 📣 calling:")

	content, err := mcpClient.CallTool(
		"fetch",
		map[string]interface{}{
			"url": "https://raw.githubusercontent.com/parakeet-nest/parakeet/refs/heads/main/blogposts/mcp-sample/demo/README.md",
		},
	)

	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println("🌍 Content:", content.Text)

}
