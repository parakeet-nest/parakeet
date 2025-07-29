package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/gear"
	"github.com/parakeet-nest/parakeet/llm"
	mcphttp "github.com/parakeet-nest/parakeet/mcp-http"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡", err)
	}

	ollamaUrl := gear.GetEnvString("OLLAMA_HOST", "http://localhost:11434")
	modelWithToolsSupport := gear.GetEnvString("LLM_WITH_TOOLS_SUPPORT", "qwen2.5:latest")
	mcpHTTPServerUrl := gear.GetEnvString("MCP_HOST", "http://0.0.0.0:9090")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new mcp HTTP client
	mcpClient, err := mcphttp.NewClient(ctx, mcpHTTPServerUrl)

	if err != nil {
		log.Fatalln("😡", err)
	}
	err = mcpClient.Start()
	if err != nil {
		log.Fatalln("😡", err)
	}
	result, err := mcpClient.Initialize()
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println("🚀 Initialized with server:", result.ServerInfo.Name, result.ServerInfo.Version)

	ollamaTools, err := mcpClient.ListTools()
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println("🛠️ Tools found:")
	for _, tool := range ollamaTools {
		fmt.Println("  -", tool.Function.Name, ":", tool.Function.Description)
	}

	// Create the chat query
	messages := []llm.Message{
		{Role: "user", Content: "Can you roll 3 dices with 6 sides?"},
	}

	options := llm.SetOptions(map[string]any{
		option.Temperature:   0.2,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.0,
	})

	toolsQuery := llm.Query{
		Model:    modelWithToolsSupport,
		Messages: messages,
		Tools:    ollamaTools,
		Options:  options,
	}

	fmt.Println("\n🤖 Chat completion...")
	chatCompletion, err := completion.Chat(ollamaUrl, toolsQuery)
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println("🟢 🦙 Answer:", chatCompletion.Message.Content)

	fmt.Println()

	if len(chatCompletion.Message.ToolCalls) > 0 {
		fmt.Println("🟢 🦙 Tool Calls:", chatCompletion.Message.ToolCalls)
		
		// Find the tool to call (assuming first tool call for simplicity)
		if len(chatCompletion.Message.ToolCalls) > 0 {
			toolCall := chatCompletion.Message.ToolCalls[0]
			fmt.Println("🔧 Calling tool:", toolCall.Function.Name)
			fmt.Println("  Arguments:", toolCall.Function.Arguments)

			result, err := mcpClient.CallTool(toolCall.Function.Name, toolCall.Function.Arguments)
			if err != nil {
				log.Fatalln("😡", err)
			}
			fmt.Println("  Result:", result.Text)
		}
	}

	defer mcpClient.Close()
}