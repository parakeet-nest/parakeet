package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
)

func TestWithOpenAI(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	openaiUrl := os.Getenv("OPENAI_BASE_URL")+"/v1"
	if openaiUrl == "" {
		openaiUrl = "https://api.openai.com/v1"
	}

	model := os.Getenv("OPENAI_LLM_CHAT")
	if model == "" {
		model = "gpt-4o-mini"
	}

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
	})

	schema := map[string]any{
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
	}

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "user", Content: "Tell me about Canada"},
		},
		Options: options,
		Format:  schema,
		Raw:     false,
	}


	answer, err := completion.Chat(openaiUrl, query, provider.OpenAI, os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		// test if the model is not found
		if modelErr, ok := err.(*completion.ModelNotFoundError); ok {
			fmt.Printf("💥 Got Model Not Found error: %s\n", modelErr.Message)
			fmt.Printf("😡 Error code: %d\n", modelErr.Code)
			fmt.Printf("🧠 Expected Model: %s\n", modelErr.Model)
		}
		if noHostErr, ok := err.(*completion.NoSuchOllamaHostError); ok {
			fmt.Printf("🦙 Got No Such Ollama Host error: %s\n", noHostErr.Message)
			fmt.Printf("🌍 Expected Host: %s\n", noHostErr.Host)
		}
		log.Fatal("😡:", err)
		t.Fatalf("Error: %v", err)
	}
	fmt.Println("OpenAI:", answer.Message.Content)
}
