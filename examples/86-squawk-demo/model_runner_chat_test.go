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

func TestChatWithModelRunner(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	modelRunnerBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")
	model := os.Getenv("MODEL_RUNNER_LLM_CHAT")

	dmrParrot := squawk.New().Model(model).BaseURL(modelRunnerBaseUrl).Provider(provider.DockerModelRunner)

	dmrParrot.
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
	fmt.Println("🐳 Answer:", dmrParrot.LastAnswer().Message.Content)
	fmt.Println("====================================")

}
