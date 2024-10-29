package main

import (
	"os"

	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/flock"
	"github.com/parakeet-nest/parakeet/ui"
	"github.com/parakeet-nest/parakeet/ui/colors"
)

func PrintMessages(messages []llm.Message) {
	for _, message := range messages {
		ui.Println(colors.Blue, "-", message.Role, ": ", message.Content)
	}
}

func main() {

	ollamaUrl := os.Getenv("OLLAMA_URL")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.7,
		option.RepeatLastN:   8,   // Increased for better context
		option.RepeatPenalty: 1.2, // Reduced further to help with response generation
		option.TopK:          40,
		option.TopP:          0.9,
	})

	bob := flock.Agent{
		Name:      "Bob",
		Model:     "qwen2.5:3b",
		OllamaUrl: ollamaUrl,
		Options:   options,
	}

	bob.SetInstructions(func(contextVars map[string]interface{}) string {

		userName, ok := contextVars["userName"].(string)

		if !ok {
			return "Help the user with whatever they want."
		}

		return fmt.Sprintf(`Help the user, and always call him by his name first: %s, 
		do whatever he want.`, userName)
	})

	orchestrator := flock.Orchestrator{}

	bobResponse, err := orchestrator.Run(
		bob,
		[]llm.Message{
			{Role: "user", Content: "Hello, What is the best pizza in the world?"},
		},
		map[string]interface{}{"userName": "Sam"},
	)

	if err != nil {
		log.Fatalln("😡:", err)
	}

	ui.Println(colors.Green, "🤖 Bob's response: ")
	ui.Println(colors.Yellow, bobResponse.Agent)
	ui.Println(colors.Blue, bobResponse.ContextVariables)
	ui.Println(colors.Purple, bobResponse.Messages)

	PrintMessages(bobResponse.Messages)

	lastBobMessage := bobResponse.Messages[len(bobResponse.Messages)-1]
	ui.Println(colors.Green, "🤖 Last Bob's message: ", lastBobMessage.Content)

}
