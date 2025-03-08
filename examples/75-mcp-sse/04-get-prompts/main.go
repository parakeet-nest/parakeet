package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/gear"
	mcpsse "github.com/parakeet-nest/parakeet/mcp-sse"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡", err)
	}

	//ollamaUrl := gear.GetEnvString("OLLAMA_HOST", "http://localhost:11434")
	//modelWithToolsSupport := gear.GetEnvString("LLM_WITH_TOOLS_SUPPORT", "qwen2.5:0.5b")
	//chatModel := gear.GetEnvString("LLM_CHAT", "qwen2.5:0.5b")
	mcpSSEServerUrl := gear.GetEnvString("MCP_HOST", "http://0.0.0.0:5001")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new mcp client
	mcpClient, err := mcpsse.NewClient(ctx, mcpSSEServerUrl)

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

	prompts, err := mcpClient.ListPrompts()
	if err != nil {
		log.Fatalln("😡", err)
	}
	fmt.Println("📦", prompts)

	fmt.Println("📣 List all Prompts:")
	for _, prompt := range prompts {
		fmt.Println("- Name:", prompt.Name, "Arguments:", prompt.Arguments)
	}

	// Find a specific prompt
	fmt.Println("🔍 Find a specific prompt:")
	promptInfo, _ := prompts.Find("summarize")
	fmt.Println("📣 Prompt:", promptInfo.Name, "Arguments:", promptInfo.Arguments)
	//prompt.Arguments[0].Name

	prompt, _ := mcpClient.GetAndFillPrompt(promptInfo.Name, map[string]string{"content": "This is the text of the content."})

	fmt.Println("📣 Filled Prompt:", "role:", prompt.Messages[0].Role, "content:", prompt.Messages[0].Content)

	mcpClient.Close()
	fmt.Println("👋 Bye!")
}
