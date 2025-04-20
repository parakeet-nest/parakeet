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

func displaySettings(ollamaUrl, modelWithToolsSupport, chatModel, mcpSSEServerUrl string) {
	fmt.Println("🦙 OLLAMA_HOST:", ollamaUrl)
	fmt.Println("🛠️ LLM_WITH_TOOLS_SUPPORT:", modelWithToolsSupport)
	fmt.Println("🤖 LLM_CHAT:", chatModel)
	fmt.Println("🔌 MCP_HOST:", mcpSSEServerUrl)
}

func main() {

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡", err)
	}

	// Get environment variables
	ollamaUrl := gear.GetEnvString("OLLAMA_HOST", "http://localhost:11434")
	modelWithToolsSupport := gear.GetEnvString("LLM_WITH_TOOLS_SUPPORT", "qwen2.5:0.5b")
	chatModel := gear.GetEnvString("LLM_CHAT", "qwen2.5:0.5b")
	mcpSSEServerUrl := gear.GetEnvString("MCP_HOST", "http://0.0.0.0:5001")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new mcp client
	mcpClient, err := mcpsse.NewClient(ctx, mcpSSEServerUrl)
	if err != nil {
		log.Fatalln("😡", err)
	}
	// Start the client
	err = mcpClient.Start()
	if err != nil {
		log.Fatalln("😡", err)
	}
	// Initialize the client
	result, err := mcpClient.Initialize()
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println("1. 🚀 Initialized with server:", result.ServerInfo.Name, result.ServerInfo.Version)
	fmt.Println("2. 📚 Reading resource from the MCP server...")

	resourceResult, err := mcpClient.ReadResource("tools-system://instructions")
	if err != nil {
		log.Fatalln("😡", err)
	}
	toolsSystemInstructions := resourceResult.Contents

	resourceResult, err = mcpClient.ReadResource("chat-system://instructions")
	if err != nil {
		log.Fatalln("😡", err)
	}
	chatSystemInstructions := resourceResult.Contents

	fmt.Println("  - 📚 Tools System Instructions:", toolsSystemInstructions)
	fmt.Println("  - 📚 Chat System Instructions:", chatSystemInstructions)

	fmt.Println("3. 📝 Get tools Prompt from the MCP server...")
	// Get the prompt for the Tools LLM from the MCP server
	promptForToolsLLM, err := mcpClient.GetAndFillPrompt(
		"fetch-page",
		map[string]string{"url": "https://raw.githubusercontent.com/sea-monkeys/WASImancer/main/README.md"},
	)
	if err != nil {
		log.Fatalln("😡", err)
	}

	fmt.Println("4. 📣 Filled Prompt:", "role:", promptForToolsLLM.Messages[0].Role, "content:", promptForToolsLLM.Messages[0].Content)

	fmt.Println("5. 🛠️ Get tools list from the MCP server...")
	
	// Get the list of tools from the MCP server
	ollamaTools, err := mcpClient.ListTools()
	if err != nil {
		log.Fatalln("😡", err)
	}

	// Prepare messages for the Tools LLM
	messagesForToolsLLM := []llm.Message{
		{Role: "system", Content: toolsSystemInstructions[0]},
	}
	messagesForToolsLLM = append(messagesForToolsLLM, promptForToolsLLM.Messages...)

	// Set options for the Tools LLM
	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
	})

	// Prepare query for the Tools LLM
	toolsQuery := llm.Query{
		Model:    modelWithToolsSupport,
		Messages: messagesForToolsLLM,
		Tools:    ollamaTools,
		Options:  options,
		Format:   "json",
	}

	fmt.Println("6. 📣 Send tools request to the LLM...")
	// Call the Tools LLM
	answer, err := completion.Chat(ollamaUrl, toolsQuery)
	if err != nil {
		log.Fatalln("😡", err)
	}

	// Search tool(s) to call for execution in the answer
	tool, err := answer.Message.ToolCalls.Find("fetch")
	if err != nil {
		log.Fatalln("😡", err)
	}
	fmt.Println("  - 🛠️ Tool to call:", tool)

	fmt.Println("7. 🛠️ Ask the MCP server to execute the fetch tool...")
	// 🖐️ Ask the MCP server to execute the tool
	pageContent, err := mcpClient.CallTool(tool.Function.Name, tool.Function.Arguments)
	if err != nil {
		log.Fatalln("😡", err)
	}
	fmt.Println("  - 🌍 Content length:", len(pageContent.Text))

	fmt.Println("8. 📝 Get chat Prompt from the MCP server...")
	prompt, _ := mcpClient.GetAndFillPrompt("summarize", map[string]string{"content": pageContent.Text})

	fmt.Println("  - 📣 Filled Prompt:", "role:", prompt.Messages[0].Role, "content length:",len(prompt.Messages[0].Content))

	// Prepare messages for the Chat LLM
	messagesForChatLLM := []llm.Message{
		{Role: "system", Content: chatSystemInstructions[0]},
	}
	messagesForChatLLM = append(messagesForChatLLM, prompt.Messages...)


	chatOptions := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 3.0,
	})

	query := llm.Query{
		Model:    chatModel,
		Messages: messagesForChatLLM,
		Options:  chatOptions,
	}


	fmt.Println("9. 📣 Send chat request to the LLM and display the summary of the page...")
	// Call the Chat LLM
	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatalln("😡", err)
	}

	mcpClient.Close()

	fmt.Println()
	fmt.Println("👋 Bye!")
}
