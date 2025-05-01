package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/squawk"
)

func TestChatWithOllama(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	ollamaBaseUrl := os.Getenv("OLLAMA_BASE_URL")
	model := os.Getenv("OLLAMA_LLM_CHAT")
	
	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.2,
	})

	ollamaParrot := squawk.New().Model(model).BaseURL(ollamaBaseUrl).Provider(provider.Ollama)

	ollamaParrot.
		Options(options).
		System("You are a useful AI agent, you are a Star Trek expert.").
		User("Who is James T Kirk?").
		Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
			if err != nil {
				fmt.Println("😡 Error:", err)
			}
			t.Log(answer.Message.Content)
		})
	
	fmt.Println("====================================")
	fmt.Println("🦙 Answer:", ollamaParrot.LastAnswer().Message.Content)
	fmt.Println("====================================")

}
