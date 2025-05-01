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

func TestWithModelRunner(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	modelRunnerBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")
	model := os.Getenv("MODEL_RUNNER_LLM_CHAT")

	dmrParrot := squawk.New().Model(model).BaseURL(modelRunnerBaseUrl).Provider(provider.DockerModelRunner)

	dmrParrot.Schema(
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type": "string",
				},
				"capital": map[string]any{
					"type": "string",
				},
				"languages": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "string",
					},
				},
			},
			"required": []string{"name", "capital", "languages"},
		}).
		User("Tell me about Canada").
		Options(llm.SetOptions(map[string]interface{}{
			option.Temperature: 0.0,
		})).
		StructuredOutput(func(answer llm.Answer, self *squawk.Squawk, err error) {
			if err != nil {
				t.Fatalf("😡 Error: %v", err)
			}
			fmt.Println("Model Runner Structured Output:\n", answer.Message.Content)
		})

}
