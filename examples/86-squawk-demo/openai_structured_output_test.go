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

func TestWithOpenAI(t *testing.T) {
	// create a `openai.key.env` file with the following content:
	// OPENAI_API_KEY=your_openai_api_key
	err := godotenv.Load(".env", "openai.key.env")	
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	
	openAIBaseUrl := os.Getenv("OPENAI_BASE_URL")
	model := os.Getenv("OPENAI_LLM_CHAT")

	openAIParrot := squawk.New().Model(model).BaseURL(openAIBaseUrl).Provider(provider.OpenAI, os.Getenv("OPENAI_API_KEY"))

	openAIParrot.Schema(
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
			fmt.Println("OpenAI Structured Output:\n", answer.Message.Content)
		})

}
