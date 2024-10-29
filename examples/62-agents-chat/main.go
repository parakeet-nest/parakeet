package main

import (
	"os"

	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"

	"github.com/parakeet-nest/parakeet/flock"
	"github.com/parakeet-nest/parakeet/ui"
	"github.com/parakeet-nest/parakeet/ui/colors"
)

func PrintConversationHistory(conversation []llm.Message) {
	for _, message := range conversation {
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
		Name: "Bob",
		//Model:     "nemotron-mini:latest",
		Model:     "qwen2.5:1.5b",
		OllamaUrl: ollamaUrl,
		Options:   options,
	}
	bob.SetInstructions(func(contextVars map[string]interface{}) string {
		return `You are Bob, having a casual chat with Sam at a coffee shop.
		Personality: friendly, curious, and engaging.

		Your role in this conversation:
		1. You are Bob - only speak as Bob
		2. Ask Sam one clear question in each response
		3. Keep responses under 50 words
		4. Never pretend to be Sam or answer for Sam
		5. React to Sam's previous response when appropriate

		Example of good response:
		"Thanks for sharing that, Sam! I love hiking too. Have you tried any of the trails in the national park?"
		`
	})

	sam := flock.Agent{
		Name: "Sam",
		//Model:     "nemotron-mini:latest",
		Model:     "qwen2.5:1.5b",
		OllamaUrl: ollamaUrl,
		Options:   options,
	}
	sam.SetInstructions(func(contextVars map[string]interface{}) string {
		return `You are Sam, chatting with Bob at a coffee shop.
		Personality: friendly, direct, and attentive.

		Your role in this conversation:
		1. You are Sam - only speak as Sam
		2. Always directly answer Bob's questions
		3. Keep responses under 50 words
		4. Add a brief follow-up question for Bob
		5. Acknowledge what Bob just said

		Example of good response:
		"Yes, I've tried the Ridge Trail! It was beautiful. What's your favorite hiking spot around here, Bob?"
		`
	})

	flockOrchestrator := flock.Orchestrator{}

	ui.Println(colors.Blue, "\n=== Starting conversation between Bob and Sam at the coffee shop ===\n")
	// Initialize conversation with a starter
	conversation := []llm.Message{
		{Role: "system", Content: "This is a conversation between Bob and Sam. They should actively engage with each other, responding to questions and comments."},
		{Role: "user", Content: "Bob and Sam meet at a coffee shop. Start a conversation about weekend plans."},
	}

	var bobResponse flock.Response
	var samResponse flock.Response

	for i := 0; i < 10; i++ {
		// Bob's turn
		bobResponse, _ = flockOrchestrator.Run(bob, conversation, map[string]interface{}{})

		lastBobMessage := bobResponse.Messages[len(bobResponse.Messages)-1]
		ui.Println(colors.Green, "🐼 Bob: "+lastBobMessage.Content)

		// Add Bob's message to conversation history
		conversation = append(conversation, llm.Message{
			Role:    "user",
			Content: lastBobMessage.Content,
		})

		// Sam's turn
		samResponse, _ = flockOrchestrator.Run(sam, conversation, map[string]interface{}{})

		lastSamMessage := samResponse.Messages[len(samResponse.Messages)-1]
		ui.Println(colors.Yellow, "🐻 Sam: "+lastSamMessage.Content)

		// Add Sam's message to conversation history
		conversation = append(conversation, llm.Message{
			Role:    "user",
			Content: lastSamMessage.Content,
		})

	}

	ui.Println(colors.Blue, "\n=== End of conversation ===\n")
}
