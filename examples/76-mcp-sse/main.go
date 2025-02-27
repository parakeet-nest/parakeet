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
	"github.com/parakeet-nest/parakeet/mcphelpers"
)

func displaySettings(ollamaUrl, modelWithToolsSupport, chatModel, mcpSSEServerUrl string) {
	fmt.Println("🦙 OLLAMA_HOST:", ollamaUrl)
	fmt.Println("🛠️ LLM_WITH_TOOLS_SUPPORT:", modelWithToolsSupport)
	fmt.Println("🤖 LLM_CHAT:", chatModel)
	fmt.Println("🔌 MCP_HOST:", mcpSSEServerUrl)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡", err)
	}

	ollamaUrl := gear.GetEnvString("OLLAMA_HOST", "http://localhost:11434")
	modelWithToolsSupport := gear.GetEnvString("LLM_WITH_TOOLS_SUPPORT", "qwen2.5:0.5b")
	chatModel := gear.GetEnvString("LLM_CHAT", "qwen2.5:0.5b")
	mcpSSEServerUrl := gear.GetEnvString("MCP_HOST", "http://http://0.0.0.0:5001")

	displaySettings(ollamaUrl, modelWithToolsSupport, chatModel, mcpSSEServerUrl)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a new mcp client
	mcpClient, result, err := mcphelpers.GetMCPSSEClient(ctx, mcpSSEServerUrl)

	if err != nil {
		log.Fatalln("😡", err)
	}
	defer mcpClient.Close()

	fmt.Println("🚀 Initialized with server:", result.ServerInfo.Name, result.ServerInfo.Version)

	ollamaTools, err := mcphelpers.GetSSETools(ctx, mcpClient)

	fmt.Println("=============================================")
	fmt.Println("🛠️ Tools:", ollamaTools)
	fmt.Println("=============================================")

	if err != nil {
		log.Fatalln("😡", err)
	}

	systemMCPInstructions := `You are a useful AI agent. 
	Your job is to understand the user prompt ans decide if you need to use a tool to run external commands.
	Ignore all things not related to the usage of a tool.
	`

	systemChatInstructions := `You are a useful AI agent. your job is to answer the user prompt.
	If you detect that the user prompt is related to a tool, ignore this part and focus on the other parts.
	`
	// This prompt will be use by the Tools LLM and the Chat LLM
	globalPrompt := `Fetch this page: https://raw.githubusercontent.com/sea-monkeys/WASImancer/main/README.md 
	and then make a brief summary of the content.`

	// Send request to a LLM with tools suppot
	messages := []llm.Message{
		{Role: "system", Content: systemMCPInstructions},
		{Role: "user", Content: globalPrompt},
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

	mcpResult, err := mcphelpers.CallSSETool(ctx, mcpClient, toolCall.Function.Name, toolCall.Function.Arguments)
	if err != nil {
		log.Fatalln("😡", err)
	}
	// Get the text from the result
	contentOfTheWebPage, _ := mcphelpers.GetTextFromResult(mcpResult)

	fmt.Println("📝 CONTENT:", contentOfTheWebPage)

	// add this {Role: "user", Content: contentForThePrompt} to the messages
	messages = append(messages,
		llm.Message{Role: "system", Content: systemChatInstructions},
		llm.Message{Role: "user", Content: globalPrompt},
		llm.Message{Role: "user", Content: contentOfTheWebPage},
	)

	chatOptions := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.5,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 3.0,
	})

	query := llm.Query{
		Model:    chatModel,
		Messages: messages,
		Options:  chatOptions,
	}

	fmt.Println()
	fmt.Println("📝 SUMMARY:")

	_, err = completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatalln("😡", err)
	}

}
