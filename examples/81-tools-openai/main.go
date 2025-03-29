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
	// create a `.env` file with the following content:
	// OPENAI_API_KEY=your_openai_api_key
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("😡:", err)
	}

	openAIURL := "https://api.openai.com/v1"
	model := "gpt-4o-mini"
	openAIKey := os.Getenv("OPENAI_API_KEY")

	toolsList := []llm.Tool{
		{
			Type: "function",
			Function: llm.Function{
				Name:        "hello",
				Description: "Say hello to a given person with his name",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"name": {
							Type:        "string",
							Description: "The name of the person",
						},
					},
					Required: []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.Function{
				Name:        "addNumbers",
				Description: "Make an addition of the two given numbers",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"a": {
							Type:        "number",
							Description: "first operand",
						},
						"b": {
							Type:        "number",
							Description: "second operand",
						},
					},
					Required: []string{"a", "b"},
				},
			},
		},
	}

	messages := []llm.Message{
		{Role: "user", Content: `say "hello" to Bob`},
	}

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
	})

	query := llm.Query{
		Model:    model,
		Messages: messages,
		Tools:    toolsList,
		Options:  options,
	}

	answer, err := completion.Chat(openAIURL, query, provider.OpenAI, openAIKey)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(answer.Message.ToolCalls)

	// Search tool to call in the answer
	tool, err := answer.Message.ToolCalls.Find("hello")
	if err != nil {
		log.Fatal("😡:", err)
	}
	result, _ := tool.Function.ToJSONString()
	fmt.Println(result)

	messages = []llm.Message{
		{Role: "user", Content: `add 2 and 40`},
	}

	query = llm.Query{
		Model:    model,
		Messages: messages,
		Tools:    toolsList,
		Options:  options,
	}

	answer, err = completion.Chat(openAIURL, query, provider.OpenAI, openAIKey)
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Search tool to call in the answer
	tool, err = answer.Message.ToolCalls.Find("addNumbers")
	if err != nil {
		log.Fatal("😡:", err)
	}
	result, _ = tool.Function.ToJSONString()
	fmt.Println(result)

}
