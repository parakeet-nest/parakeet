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

func TestChatWithOpenAI(t *testing.T) {
	err := godotenv.Load(".env", "openai.key.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	openAIBaseUrl := os.Getenv("OPENAI_BASE_URL")
	model := os.Getenv("OPENAI_LLM_CHAT")

	openAIParrot := squawk.New().
		Model(model).
		BaseURL(openAIBaseUrl).
		Provider(provider.OpenAI, os.Getenv("OPENAI_API_KEY"))

	openAIParrot.
		Options(llm.SetOptions(map[string]interface{}{
			option.Temperature:   0.0,
			option.RepeatLastN:   2,
			option.RepeatPenalty: 2.2,
		})).
		System("You are a useful AI agent, you are a Star Trek expert.").
		User("Who is James T Kirk?").
		Chat(func(answer llm.Answer, self *squawk.Squawk, err error) {
			if err != nil {
				fmt.Println("😡 Error:", err)
			}
			t.Log(answer.Message.Content)
		})

	fmt.Println("====================================")
	fmt.Println("⚫️ Answer:", openAIParrot.LastAnswer().Message.Content)
	fmt.Println("====================================")

}
