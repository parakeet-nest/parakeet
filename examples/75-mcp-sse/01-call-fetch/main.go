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
	mcpsse "github.com/parakeet-nest/parakeet/mcp-sse"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡", err)
	}

	ollamaUrl := gear.GetEnvString("OLLAMA_HOST", "http://localhost:11434")
	modelWithToolsSupport := gear.GetEnvString("LLM_WITH_TOOLS_SUPPORT", "qwen2.5:0.5b")
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

	ollamaTools, err := mcpClient.ListTools()

	fmt.Println("=============================================")
	fmt.Println("🛠️ Tools:", ollamaTools)
	fmt.Println("=============================================")

	if err != nil {
		log.Fatalln("😡", err)
	}
	toolPrompt := `Fetch this page: https://raw.githubusercontent.com/sea-monkeys/WASImancer/main/README.md
	`
	// Send request to a LLM with tools suppot
	messages := []llm.Message{
		{Role: "user", Content: toolPrompt},
	}

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.0,
	})

	toolsQuery := llm.Query{
		Model:    modelWithToolsSupport,
		Messages: messages,
		Tools:    ollamaTools,
		Options:  options,
		Format:   "json",
	}

	answer, err := completion.Chat(ollamaUrl, toolsQuery)
	if err != nil {
		log.Fatalln("😡", err)
	}

	// Search tool to call in the answer
	tool, err := answer.Message.ToolCalls.Find("fetch")
	if err != nil {
		log.Fatalln("😡", err)
	}
	fmt.Println("🛠️ Tool to call:", tool)

	// 🖐️ Call the mcp server
	fmt.Println("🦙🛠️ 📣 calling:", tool.Function.Name, tool.Function.Arguments)

	content, err := mcpClient.CallTool(tool.Function.Name, tool.Function.Arguments)
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println("🌍 Content:", content.Text)

	mcpClient.Close()
	fmt.Println("👋 Bye!")
}
