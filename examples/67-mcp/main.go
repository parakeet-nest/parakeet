package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"

	mcpstdio "github.com/parakeet-nest/parakeet/mcp-stdio"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡", err)
	}

	ollamaUrl := os.Getenv("OLLAMA_HOST")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	modelWithToolsSupport := os.Getenv("LLM_WITH_TOOLS_SUPPORT")
	if modelWithToolsSupport == "" {
		modelWithToolsSupport = "qwen2.5:0.5b"
	}

	chatModel := os.Getenv("LLM_CHAT")
	if chatModel == "" {
		chatModel = "qwen2.5:0.5b"
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new mcp client
	mcpClient, err := mcpstdio.NewClient(ctx, "docker",
		[]string{}, // Empty ENV
		"run",
		"--rm",
		"-i",
		"mcp-curl",
	)

	if err != nil {
		log.Fatalln("😡", err)
	}
	//defer mcpClient.Close()

	_, err = mcpClient.Initialize()
	if err != nil {
		log.Fatalln("😡", err)
	}

	ollamaTools, err := mcpClient.ListTools()

	if err != nil {
		log.Fatalln("😡", err)
	}

	// Send request to a LLM with tools suppot
	messages := []llm.Message{
		{Role: "user", Content: "Fetch this page: https://raw.githubusercontent.com/parakeet-nest/parakeet/main/README.md"},
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

	// Get the first tool call from the answer (what the LLM wants to do / understand)
	toolCall := answer.Message.ToolCalls[0]

	// 🖐️ Call the mcp server
	fmt.Println("🦙🛠️ 📣 calling:", toolCall.Function.Name, toolCall.Function.Arguments)

	mcpResult, err := mcpClient.CallTool(toolCall.Function.Name, toolCall.Function.Arguments)
	if err != nil {
		log.Fatalln("😡", err)
	}
	// Get the text from the result
	contentOfTheWebPage := mcpResult.Text

	//fmt.Println("🌍 Content:", contentOfTheWebPage)

	// add this {Role: "user", Content: contentForThePrompt} to the messages
	messages = append(messages,
		llm.Message{Role: "user", Content: "MWhat is the main topic of the below page:"},
		llm.Message{Role: "user", Content: contentOfTheWebPage},
	)

	chatOptions := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 3.0,
		option.NumCtx:        content.EstimateGPTTokens(contentOfTheWebPage) + 1000,
	})

	query := llm.Query{
		Model:    chatModel,
		Messages: messages,
		Options:  chatOptions,
	}

	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatalln("😡", err)
	}

	mcpClient.Close()
}
