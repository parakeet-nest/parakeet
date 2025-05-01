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

func TestWithOllama(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	
	ollamaBaseUrl := os.Getenv("OLLAMA_BASE_URL")
	model := os.Getenv("OLLAMA_LLM_CHAT")

	ollamaParrot := squawk.New().Model(model).BaseURL(ollamaBaseUrl).Provider(provider.Ollama)

	ollamaParrot.Schema(
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
			fmt.Println("Ollama Structured Output:\n", answer.Message.Content)
		})


}
