/*
Topic: Parakeet
Generate a chat completion with Ollama and parakeet
no streaming
*/

package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	ollamaUrl := os.Getenv("OLLAMA_HOST")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	model := os.Getenv("LLM_CHAT")
	if model == "" {
		model = "qwen2.5:0.5b"
	}

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
	})

	// define schema for a structured output
	// ref: https://ollama.com/blog/structured-outputs
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

	//model = "ai/qwen2.5:latest"

	query := llm.Query{
		Model: model,
		Messages: []llm.Message{
			{Role: "user", Content: "Tell me about Canada"},
		},
		Options: options,
		Format:  schema,
		Raw:     false,
	}



	// ollamaUrl = "http://localhost:12434/engines/llama.cpp/v1"
	// TODO: if provider is DMR or OpenAI , Copy the value of Fromat to ResponseFormat (response_format)

	//answer, err := completion.Chat(ollamaUrl, query, provider.DockerModelRunner)

	answer, err := completion.Chat(ollamaUrl, query, provider.Ollama)
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

	}
	fmt.Println(answer.Message.Content)

}
